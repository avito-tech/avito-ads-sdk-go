package avito

import (
	"fmt"
	"os"
	"strconv"
)

// Client — клиент API Авито Реклама.
//
// Создаётся через NewClient и предоставляет доступ к сервисам (группам методов).
// Все запросы выполняются от лица аккаунта из конфигурации (токен привязан к
// одному аккаунту). Методы сервисов принимают context.Context.
//
// Пример:
//
//	client, err := avito.NewClient("CLIENT_ID", "CLIENT_SECRET", 123456789, avito.WithSandbox())
//	if err != nil { /* ... */ }
//	balance, err := client.Account.GetBalance(ctx)
type Client struct {
	cfg *config
	tr  *transport

	Account      *AccountService
	ChildAccount *ChildAccountService
	Advertisers  *AdvertiserService
	Contracts    *ContractService
	Campaigns    *CampaignService
	Groups       *GroupService
	Creatives    *CreativeService
	Statistics   *StatisticsService
	Users        *UserService
}

// NewClient создаёт клиент из учётных данных и опций.
func NewClient(clientID, clientSecret string, accountID int64, opts ...Option) (*Client, error) {
	cfg := defaultConfig(clientID, clientSecret, accountID)
	for _, opt := range opts {
		opt(cfg)
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	tr := &transport{cfg: cfg}
	c := &Client{cfg: cfg, tr: tr}
	base := baseService{tr: tr, accountID: accountID}

	c.Account = &AccountService{base}
	c.ChildAccount = &ChildAccountService{base}
	c.Advertisers = &AdvertiserService{base}
	c.Contracts = &ContractService{base}
	c.Campaigns = &CampaignService{base}
	c.Groups = &GroupService{base}
	c.Creatives = &CreativeService{base}
	c.Statistics = &StatisticsService{base}
	c.Users = &UserService{base}

	return c, nil
}

// NewFromEnv создаёт клиент из переменных окружения.
//
// Обязательные: AVITO_ADS_CLIENT_ID, AVITO_ADS_CLIENT_SECRET, AVITO_ADS_ACCOUNT_ID.
// Необязательные: AVITO_ADS_ENVIRONMENT (production|sandbox), AVITO_ADS_MAX_RETRIES,
// AVITO_ADS_TIMEOUT_SECONDS, AVITO_ADS_TOKEN_LEEWAY_SECONDS.
//
// Дополнительные опции (extra) применяются после опций из окружения.
func NewFromEnv(extra ...Option) (*Client, error) {
	clientID := os.Getenv("AVITO_ADS_CLIENT_ID")
	clientSecret := os.Getenv("AVITO_ADS_CLIENT_SECRET")
	accountRaw := os.Getenv("AVITO_ADS_ACCOUNT_ID")
	accountID, err := strconv.ParseInt(accountRaw, 10, 64)
	if err != nil {
		return nil, configErrorf("AVITO_ADS_ACCOUNT_ID должно быть числом, получено %q", accountRaw)
	}
	opts := append(optionsFromEnv(), extra...)
	return NewClient(clientID, clientSecret, accountID, opts...)
}

// AccountID возвращает идентификатор аккаунта клиента.
func (c *Client) AccountID() int64 { return c.cfg.accountID }

// BaseURL возвращает базовый адрес API с учётом окружения.
func (c *Client) BaseURL() string { return c.cfg.baseURI() }

// baseService — общая основа всех сервисов.
type baseService struct {
	tr        *transport
	accountID int64
}

// accountPath строит путь вида v1/account/{accountID} + необязательный суффикс.
func (s baseService) accountPath(suffix string) string {
	p := fmt.Sprintf("v1/account/%d", s.accountID)
	if suffix != "" {
		p += "/" + suffix
	}
	return p
}
