package avito

import (
	"context"
	"net/http"
	"testing"
)

const balancePath = "/ads/v1/account/123456789/balance"

func TestSuccessParsesBodyAndPointBalance(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("test-token"))
	ts.on("GET", balancePath, jsonHandler(200,
		map[string]any{"balance": 15000, "bonusBalance": 500},
		map[string]string{"Api-Point-Balance": "4231"}))

	c := newTestClient(t, ts)
	bal, err := c.Account.GetBalance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if bal.Balance != 15000 || bal.BonusBalance != 500 {
		t.Errorf("balance = %+v", bal)
	}
	last := ts.requests[len(ts.requests)-1]
	if last.Auth != "Bearer test-token" {
		t.Errorf("auth = %q", last.Auth)
	}
}

func TestErrorStatusMapsToTypedError(t *testing.T) {
	cases := []struct {
		status int
		check  func(error) bool
	}{
		{400, IsBadRequest},
		{403, IsAccessDenied},
		{404, IsNotFound},
		{500, IsServerError},
	}
	for _, tc := range cases {
		ts := newTestServer()
		ts.on("POST", "/token", tokenHandler("t"))
		ts.on("GET", balancePath, jsonHandler(tc.status,
			map[string]any{"code": "e", "message": "msg"}, nil))
		c, _ := NewClient("cid", "secret", 123456789,
			WithHTTPClient(&http.Client{Transport: &rewriteTransport{target: ts.srv.URL}}),
			WithRetryBaseWait(0), WithMaxRetries(0))
		_, err := c.Account.GetBalance(context.Background())
		if err == nil || !tc.check(err) {
			t.Errorf("status %d: unexpected err %v", tc.status, err)
		}
		ts.close()
	}
}

func TestErrorCarriesCodeAndMessage(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("GET", balancePath, jsonHandler(400,
		map[string]any{"code": "bad_request", "message": "Плохой запрос"}, nil))
	c := newTestClient(t, ts)
	_, err := c.Account.GetBalance(context.Background())
	var apiErr *APIError
	if !asAPIError(err, &apiErr) {
		t.Fatalf("ожидался *APIError, получено %v", err)
	}
	if apiErr.Code != "bad_request" || apiErr.StatusCode != 400 || apiErr.Message != "Плохой запрос" {
		t.Errorf("apiErr = %+v", apiErr)
	}
}

func TestRateLimitRetryAfter(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("GET", balancePath, jsonHandler(429,
		map[string]any{"message": "too many"}, map[string]string{"Retry-After": "7"}))
	c, _ := NewClient("cid", "secret", 123456789,
		WithHTTPClient(&http.Client{Transport: &rewriteTransport{target: ts.srv.URL}}),
		WithRetryBaseWait(0), WithMaxRetries(0))
	_, err := c.Account.GetBalance(context.Background())
	var apiErr *APIError
	if !asAPIError(err, &apiErr) || !IsRateLimit(err) {
		t.Fatalf("ожидался RateLimit, получено %v", err)
	}
	if apiErr.RetryAfter != 7 {
		t.Errorf("retryAfter = %d", apiErr.RetryAfter)
	}
}

func TestRefreshesTokenAndRetriesOn401(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	// Первый токен, затем 401, затем второй токен, затем успех.
	ts.enqueue("POST", "/token", tokenHandler("tok-1"))
	ts.enqueue("GET", balancePath, jsonHandler(401, map[string]any{"message": "expired"}, nil))
	ts.enqueue("POST", "/token", tokenHandler("tok-2"))
	ts.enqueue("GET", balancePath, jsonHandler(200, map[string]any{"balance": 777}, nil))

	c := newTestClient(t, ts)
	bal, err := c.Account.GetBalance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if bal.Balance != 777 {
		t.Errorf("balance = %d", bal.Balance)
	}
	last := ts.requests[len(ts.requests)-1]
	if last.Auth != "Bearer tok-2" {
		t.Errorf("последний запрос должен идти с обновлённым токеном, получено %q", last.Auth)
	}
}

func TestPersistent401ReturnsAuthError(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.enqueue("POST", "/token", tokenHandler("tok-1"))
	ts.enqueue("GET", balancePath, jsonHandler(401, map[string]any{"message": "no"}, nil))
	ts.enqueue("POST", "/token", tokenHandler("tok-2"))
	ts.enqueue("GET", balancePath, jsonHandler(401, map[string]any{"message": "no"}, nil))
	c := newTestClient(t, ts)
	_, err := c.Account.GetBalance(context.Background())
	if !IsAuthentication(err) {
		t.Errorf("ожидался AuthenticationError, получено %v", err)
	}
}

func TestRetriesOn429ThenSucceeds(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.enqueue("GET", balancePath, jsonHandler(429, map[string]any{"message": "slow"}, nil))
	ts.enqueue("GET", balancePath, jsonHandler(200, map[string]any{"balance": 10}, nil))
	c := newTestClient(t, ts)
	bal, err := c.Account.GetBalance(context.Background())
	if err != nil || bal.Balance != 10 {
		t.Errorf("bal=%v err=%v", bal, err)
	}
}

func TestRetriesOn503ThenSucceeds(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.enqueue("GET", balancePath, jsonHandler(503, map[string]any{"message": "down"}, nil))
	ts.enqueue("GET", balancePath, jsonHandler(200, map[string]any{"balance": 11}, nil))
	c := newTestClient(t, ts)
	bal, err := c.Account.GetBalance(context.Background())
	if err != nil || bal.Balance != 11 {
		t.Errorf("bal=%v err=%v", bal, err)
	}
}

// asAPIError — обёртка над errors.As для тестов.
func asAPIError(err error, target **APIError) bool {
	for err != nil {
		if e, ok := err.(*APIError); ok {
			*target = e
			return true
		}
		type unwrapper interface{ Unwrap() error }
		if u, ok := err.(unwrapper); ok {
			err = u.Unwrap()
		} else {
			return false
		}
	}
	return false
}
