package avito

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// transport инкапсулирует выполнение запросов: авторизацию, повторы и разбор.
type transport struct {
	cfg *config
}

// tokenResponse — тело ответа эндпоинта /token.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// apiErrorBody — тело ошибки API.
type apiErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// getToken возвращает действующий токен, обновляя его при необходимости.
func (t *transport) getToken(ctx context.Context) (accessToken, error) {
	key := t.cfg.cacheKey()
	if value, expires, ok := t.cfg.tokenStorage.Get(key); ok {
		tok := accessToken{Value: value, ExpiresAt: expires, TokenType: "Bearer"}
		if !tok.isExpired(t.cfg.tokenLeeway, time.Now()) {
			return tok, nil
		}
	}
	return t.fetchToken(ctx)
}

// forgetToken сбрасывает закэшированный токен.
func (t *transport) forgetToken() {
	t.cfg.tokenStorage.Delete(t.cfg.cacheKey())
}

// fetchToken запрашивает новый токен у эндпоинта /token.
func (t *transport) fetchToken(ctx context.Context) (accessToken, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", t.cfg.clientID)
	form.Set("client_secret", t.cfg.clientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.cfg.tokenURL(),
		strings.NewReader(form.Encode()))
	if err != nil {
		return accessToken{}, fmt.Errorf("avito: не удалось создать запрос токена: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", t.cfg.userAgent)

	resp, err := t.cfg.httpClient.Do(req)
	if err != nil {
		return accessToken{}, fmt.Errorf("avito: сетевая ошибка при получении токена: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return accessToken{}, &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("не удалось получить токен (HTTP %d)", resp.StatusCode),
		}
	}

	var tr tokenResponse
	if err := json.Unmarshal(body, &tr); err != nil || tr.AccessToken == "" {
		return accessToken{}, &APIError{
			StatusCode: resp.StatusCode,
			Message:    "ответ эндпоинта /token не содержит access_token",
		}
	}

	expiresIn := tr.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 86400
	}
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
	t.cfg.tokenStorage.Set(t.cfg.cacheKey(), tr.AccessToken, expiresAt)

	return accessToken{Value: tr.AccessToken, ExpiresAt: expiresAt, TokenType: tr.TokenType}, nil
}

// do выполняет запрос к API и разбирает ответ в out (если out != nil).
// Возвращает остаток лимита из заголовка Api-Point-Balance (или -1, если нет).
func (t *transport) do(ctx context.Context, method, path string, body any, out any) (int, error) {
	fullURL := t.cfg.baseURI() + strings.TrimPrefix(path, "/")

	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return -1, fmt.Errorf("avito: не удалось сериализовать тело запроса: %w", err)
		}
	}

	attempt := 0
	tokenRefreshed := false

	for {
		tok, err := t.getToken(ctx)
		if err != nil {
			return -1, err
		}

		var reqBody io.Reader
		if payload != nil {
			reqBody = bytes.NewReader(payload)
		}
		req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
		if err != nil {
			return -1, fmt.Errorf("avito: не удалось создать запрос: %w", err)
		}
		req.Header.Set("Authorization", tok.authorizationHeader())
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", t.cfg.userAgent)
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := t.cfg.httpClient.Do(req)
		if err != nil {
			// Сетевая ошибка — повторяем в пределах бюджета.
			if attempt < t.cfg.maxRetries {
				attempt++
				if werr := t.wait(ctx, attempt, 0); werr != nil {
					return -1, werr
				}
				continue
			}
			return -1, fmt.Errorf("avito: сетевая ошибка при запросе %s %s: %w", method, path, err)
		}

		status := resp.StatusCode
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// 401 — однократно обновляем токен и повторяем без задержки.
		if status == http.StatusUnauthorized && !tokenRefreshed {
			tokenRefreshed = true
			t.forgetToken()
			continue
		}

		// 429 и 5xx — повторяем в пределах бюджета попыток.
		if (status == http.StatusTooManyRequests || status >= 500) && attempt < t.cfg.maxRetries {
			attempt++
			if werr := t.wait(ctx, attempt, retryAfterSeconds(resp)); werr != nil {
				return -1, werr
			}
			continue
		}

		if status >= 400 {
			return -1, buildAPIError(status, respBody, resp)
		}

		balance := pointBalance(resp)
		if out != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, out); err != nil {
				return balance, fmt.Errorf("avito: не удалось разобрать ответ: %w", err)
			}
		}
		return balance, nil
	}
}

// wait делает паузу перед повтором, уважая отмену контекста.
// При retryBaseWait == 0 (например, в тестах) задержка пропускается.
func (t *transport) wait(ctx context.Context, attempt int, retryAfter int) error {
	var d time.Duration
	if retryAfter > 0 {
		d = time.Duration(retryAfter) * time.Second
	} else {
		d = t.cfg.retryBaseWait * (1 << (attempt - 1))
	}
	if d <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func retryAfterSeconds(resp *http.Response) int {
	v := resp.Header.Get("Retry-After")
	if v == "" {
		return 0
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return 0
	}
	return n
}

func pointBalance(resp *http.Response) int {
	v := resp.Header.Get("Api-Point-Balance")
	if v == "" {
		return -1
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return -1
	}
	return n
}

func buildAPIError(status int, body []byte, resp *http.Response) *APIError {
	var eb apiErrorBody
	_ = json.Unmarshal(body, &eb)
	msg := eb.Message
	if msg == "" {
		msg = fmt.Sprintf("ошибка API (HTTP %d)", status)
	}
	return &APIError{
		StatusCode: status,
		Code:       eb.Code,
		Message:    msg,
		RetryAfter: retryAfterSeconds(resp),
	}
}
