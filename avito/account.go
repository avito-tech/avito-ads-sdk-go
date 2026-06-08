package avito

import "context"

// AccountService — методы работы с аккаунтом.
type AccountService struct{ baseService }

// Get возвращает реквизиты текущего аккаунта.
//
// GET v1/account/{accountID}
func (s *AccountService) Get(ctx context.Context) (*Account, error) {
	// Ответ может приходить как {"account": {...}} либо плоским объектом.
	var envelope struct {
		Account *Account `json:"account"`
	}
	if _, err := s.tr.do(ctx, "GET", s.accountPath(""), nil, &envelope); err != nil {
		return nil, err
	}
	if envelope.Account != nil {
		return envelope.Account, nil
	}
	var acc Account
	if _, err := s.tr.do(ctx, "GET", s.accountPath(""), nil, &acc); err != nil {
		return nil, err
	}
	return &acc, nil
}

// GetBalance возвращает баланс текущего аккаунта.
//
// GET v1/account/{accountID}/balance
func (s *AccountService) GetBalance(ctx context.Context) (*Balance, error) {
	var b Balance
	if _, err := s.tr.do(ctx, "GET", s.accountPath("balance"), nil, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// SandboxAccountInput — параметры создания тестового аккаунта.
type SandboxAccountInput struct {
	INN           string         `json:"inn"`
	ShortName     string         `json:"shortName"`
	LongName      string         `json:"longName"`
	OGRN          string         `json:"ogrn"`
	LegalAddress  string         `json:"legalAddress"`
	ActualAddress string         `json:"actualAddress"`
	Contact       map[string]any `json:"contact"`
	KPP           string         `json:"kpp,omitempty"`
	LegalType     string         `json:"legalType,omitempty"`
}

// CreateSandboxAccount создаёт тестовый аккаунт (только в песочнице).
//
// POST v1/account/{accountID}
func (s *AccountService) CreateSandboxAccount(ctx context.Context, in SandboxAccountInput) (*CreatedAccount, error) {
	if len(in.Contact) == 0 {
		return nil, validationErrorf("для создания аккаунта необходимо передать contact")
	}
	var out CreatedAccount
	if _, err := s.tr.do(ctx, "POST", s.accountPath(""), in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
