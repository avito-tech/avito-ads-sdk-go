package avito

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	apiHost          = "https://api.avito.ru"
	productionPrefix = "ads"
	sandboxPrefix    = "ads-sandbox"
)

// config хранит параметры подключения. Заполняется через функциональные опции.
type config struct {
	clientID     string
	clientSecret string
	accountID    int64
	pathPrefix   string

	httpClient    *http.Client
	maxRetries    int
	retryBaseWait time.Duration
	tokenLeeway   time.Duration
	tokenStorage  TokenStorage
	userAgent     string
}

// Option — функциональная опция конфигурации клиента.
type Option func(*config)

// WithSandbox переключает клиент на окружение песочницы.
func WithSandbox() Option {
	return func(c *config) { c.pathPrefix = sandboxPrefix }
}

// WithProduction переключает клиент на боевое окружение (используется по умолчанию).
func WithProduction() Option {
	return func(c *config) { c.pathPrefix = productionPrefix }
}

// WithHTTPClient задаёт собственный *http.Client (таймауты, транспорт, прокси).
func WithHTTPClient(client *http.Client) Option {
	return func(c *config) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithMaxRetries задаёт максимальное число повторов при 429/5xx/сетевых ошибках.
func WithMaxRetries(n int) Option {
	return func(c *config) {
		if n >= 0 {
			c.maxRetries = n
		}
	}
}

// WithRetryBaseWait задаёт базовую задержку повтора (экспоненциальный backoff).
func WithRetryBaseWait(d time.Duration) Option {
	return func(c *config) {
		if d >= 0 {
			c.retryBaseWait = d
		}
	}
}

// WithTokenLeeway задаёт запас на досрочное обновление токена.
func WithTokenLeeway(d time.Duration) Option {
	return func(c *config) {
		if d >= 0 {
			c.tokenLeeway = d
		}
	}
}

// WithTokenStorage задаёт собственное хранилище токена.
func WithTokenStorage(storage TokenStorage) Option {
	return func(c *config) {
		if storage != nil {
			c.tokenStorage = storage
		}
	}
}

// WithUserAgent задаёт значение заголовка User-Agent.
func WithUserAgent(ua string) Option {
	return func(c *config) { c.userAgent = ua }
}

// defaultConfig возвращает конфигурацию со значениями по умолчанию.
func defaultConfig(clientID, clientSecret string, accountID int64) *config {
	return &config{
		clientID:      clientID,
		clientSecret:  clientSecret,
		accountID:     accountID,
		pathPrefix:    productionPrefix,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		maxRetries:    4,
		retryBaseWait: time.Second,
		tokenLeeway:   60 * time.Second,
		tokenStorage:  NewMemoryTokenStorage(),
		userAgent:     "avito-ads-sdk-go/0.1.0",
	}
}

// validate проверяет обязательные поля конфигурации.
func (c *config) validate() error {
	if c.clientID == "" {
		return configErrorf("не задан clientID")
	}
	if c.clientSecret == "" {
		return configErrorf("не задан clientSecret")
	}
	if c.accountID <= 0 {
		return configErrorf("accountID должен быть положительным, передано %d", c.accountID)
	}
	return nil
}

// baseURI возвращает базовый адрес API с учётом окружения.
func (c *config) baseURI() string {
	return apiHost + "/" + c.pathPrefix + "/"
}

// tokenURL возвращает адрес эндпоинта получения токена (общий для prod и sandbox).
func (c *config) tokenURL() string {
	return apiHost + "/token"
}

// cacheKey — стабильный ключ кэширования токена, не раскрывающий секрет.
func (c *config) cacheKey() string {
	sum := sha1.Sum([]byte(c.clientID + ":" + c.clientSecret))
	return "avito_ads_token_" + hex.EncodeToString(sum[:])
}

// optionsFromEnv собирает опции из переменных окружения с префиксом AVITO_ADS_.
//
// Используется в NewFromEnv. Распознаёт:
//
//	AVITO_ADS_ENVIRONMENT (production|sandbox)
//	AVITO_ADS_MAX_RETRIES, AVITO_ADS_TIMEOUT_SECONDS, AVITO_ADS_TOKEN_LEEWAY_SECONDS
func optionsFromEnv() []Option {
	var opts []Option
	if os.Getenv("AVITO_ADS_ENVIRONMENT") == "sandbox" {
		opts = append(opts, WithSandbox())
	}
	if v := os.Getenv("AVITO_ADS_MAX_RETRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			opts = append(opts, WithMaxRetries(n))
		}
	}
	if v := os.Getenv("AVITO_ADS_TIMEOUT_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			opts = append(opts, WithHTTPClient(&http.Client{Timeout: time.Duration(n) * time.Second}))
		}
	}
	if v := os.Getenv("AVITO_ADS_TOKEN_LEEWAY_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			opts = append(opts, WithTokenLeeway(time.Duration(n)*time.Second))
		}
	}
	return opts
}
