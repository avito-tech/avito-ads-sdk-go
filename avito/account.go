package avito

import (
	"context"
	"errors"
)

// AccountService — методы работы с аккаунтом.
type AccountService struct{ baseService }

// Get возвращает реквизиты текущего аккаунта.
//
// GET v1/account/{accountID}
func (s *AccountService) Get(ctx context.Context) (*Account, error) {
	// Ответ всегда приходит в обёртке {"account": {...}} (V1GetAccountByIdOut).
	var envelope struct {
		Account *Account `json:"account"`
	}
	if _, err := s.tr.do(ctx, "GET", s.accountPath(""), nil, &envelope); err != nil {
		return nil, err
	}
	if envelope.Account == nil {
		return nil, errAccountMissing
	}
	return envelope.Account, nil
}

// errAccountMissing — ответ GET v1/account/{accountID} не содержит поле account.
var errAccountMissing = errors.New("avito: ответ не содержит поле account")

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
//
// Обязательные поля (swagger V1CreateAccountIn): INN, ShortName, LongName,
// OGRN, LegalAddress, LegalType (LegalTypeUL или LegalTypeIP) и Contact
// (name, phone). ActualAddress и KPP необязательны: без ActualAddress
// используется LegalAddress.
type SandboxAccountInput struct {
	INN           string         `json:"inn"`
	ShortName     string         `json:"shortName"`
	LongName      string         `json:"longName"`
	OGRN          string         `json:"ogrn"`
	LegalAddress  string         `json:"legalAddress"`
	ActualAddress string         `json:"actualAddress,omitempty"`
	Contact       map[string]any `json:"contact"`
	KPP           string         `json:"kpp,omitempty"`
	LegalType     string         `json:"legalType"`
}

// CreateSandboxAccount создаёт тестовый аккаунт (только в песочнице).
//
// Ответ 207 тоже считается успехом: аккаунт создан, но часть подготовки в
// песочнице не выполнена. Описание того, что не выполнено, — в
// CreatedAccount.Warnings (пустой список — аккаунт подготовлен полностью).
//
// POST v1/account/{accountID}
func (s *AccountService) CreateSandboxAccount(ctx context.Context, in SandboxAccountInput) (*CreatedAccount, error) {
	if len(in.Contact) == 0 {
		return nil, validationErrorf("для создания аккаунта необходимо передать contact")
	}
	if in.LegalType == "" {
		return nil, validationErrorf("для создания аккаунта необходимо передать legalType (%s или %s)", LegalTypeUL, LegalTypeIP)
	}
	var out CreatedAccount
	if _, err := s.tr.do(ctx, "POST", s.accountPath(""), in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
