package avito

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// recordedRequest хранит сведения о принятом запросе для проверок в тестах.
type recordedRequest struct {
	Method string
	Path   string
	Auth   string
	UA     string
	Body   string
}

// testServer — обёртка над httptest.Server с маршрутизацией и записью запросов.
type testServer struct {
	srv      *httptest.Server
	requests []recordedRequest
	// handlers сопоставляет "METHOD /path" с функцией-обработчиком.
	handlers map[string]http.HandlerFunc
	// queue — последовательные ответы для одного и того же маршрута (для повторов).
	queue map[string][]http.HandlerFunc
}

func newTestServer() *testServer {
	ts := &testServer{
		handlers: map[string]http.HandlerFunc{},
		queue:    map[string][]http.HandlerFunc{},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", ts.dispatch)
	ts.srv = httptest.NewServer(mux)
	return ts
}

func (ts *testServer) close() { ts.srv.Close() }

func (ts *testServer) dispatch(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	ts.requests = append(ts.requests, recordedRequest{
		Method: r.Method,
		Path:   r.URL.Path,
		Auth:   r.Header.Get("Authorization"),
		UA:     r.Header.Get("User-Agent"),
		Body:   string(bodyBytes),
	})
	// Восстанавливаем тело на случай, если обработчику оно нужно.
	r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	key := r.Method + " " + r.URL.Path
	if q, ok := ts.queue[key]; ok && len(q) > 0 {
		h := q[0]
		ts.queue[key] = q[1:]
		h(w, r)
		return
	}
	if h, ok := ts.handlers[key]; ok {
		h(w, r)
		return
	}
	http.Error(w, `{"code":"not_found","message":"нет обработчика"}`, http.StatusNotFound)
}

// on регистрирует постоянный обработчик маршрута.
func (ts *testServer) on(method, path string, h http.HandlerFunc) {
	ts.handlers[method+" "+path] = h
}

// enqueue добавляет последовательный ответ для маршрута (для проверки повторов).
func (ts *testServer) enqueue(method, path string, h http.HandlerFunc) {
	key := method + " " + path
	ts.queue[key] = append(ts.queue[key], h)
}

// jsonHandler возвращает обработчик, отдающий JSON со статусом и заголовками.
func jsonHandler(status int, body any, headers map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		switch b := body.(type) {
		case string:
			_, _ = io.WriteString(w, b)
		default:
			_ = json.NewEncoder(w).Encode(b)
		}
	}
}

// tokenHandler отдаёт валидный ответ эндпоинта /token.
func tokenHandler(token string) http.HandlerFunc {
	return jsonHandler(http.StatusOK, map[string]any{
		"access_token": token,
		"expires_in":   86400,
		"token_type":   "Bearer",
	}, nil)
}

// newTestClient создаёт клиент, направленный на тестовый сервер, с нулевой задержкой повторов.
//
// Базовый URL клиента подменяется на адрес тестового сервера через кастомный
// http.Client с RoundTripper, переписывающим хост.
func newTestClient(t *testing.T, ts *testServer) *Client {
	t.Helper()
	hc := &http.Client{Transport: &rewriteTransport{target: ts.srv.URL}}
	c, err := NewClient("cid", "secret", 123456789,
		WithHTTPClient(hc),
		WithRetryBaseWait(0),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

// rewriteTransport перенаправляет все запросы на target, сохраняя путь и query.
type rewriteTransport struct {
	target string
}

func (rt *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Переписываем схему и хост на адрес тестового сервера.
	u := *req.URL
	parts := strings.SplitN(strings.TrimPrefix(rt.target, "http://"), "/", 2)
	u.Scheme = "http"
	u.Host = parts[0]
	req2 := req.Clone(req.Context())
	req2.URL = &u
	req2.Host = u.Host
	return http.DefaultTransport.RoundTrip(req2)
}
