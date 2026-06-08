package avito

import "context"

// ChildAccountService — методы работы с дочерними аккаунтами и переводами.
type ChildAccountService struct{ baseService }

// List возвращает список дочерних аккаунтов.
//
// GET v1/account/{accountID}/children
func (s *ChildAccountService) List(ctx context.Context) ([]ChildAccount, error) {
	return s.listChildren(ctx, s.accountPath("children"))
}

// ListWithBalances возвращает список дочерних аккаунтов с балансами.
//
// GET v1/account/{accountID}/children-with-balances
func (s *ChildAccountService) ListWithBalances(ctx context.Context) ([]ChildAccount, error) {
	return s.listChildren(ctx, s.accountPath("children-with-balances"))
}

func (s *ChildAccountService) listChildren(ctx context.Context, path string) ([]ChildAccount, error) {
	var out struct {
		Children []ChildAccount `json:"children"`
	}
	if _, err := s.tr.do(ctx, "GET", path, nil, &out); err != nil {
		return nil, err
	}
	return out.Children, nil
}

// CreateNonpayerChild создаёт дочерний аккаунт-непокупатель.
//
// POST v1/account/{accountID}/create-nonpayer-child-account
func (s *ChildAccountService) CreateNonpayerChild(ctx context.Context, shortName string, selfAdvertisingEnabled bool) (*CreatedChildAccount, error) {
	body := map[string]any{
		"shortName":                shortName,
		"isSelfAdvertisingEnabled": selfAdvertisingEnabled,
	}
	var out CreatedChildAccount
	if _, err := s.tr.do(ctx, "POST", s.accountPath("create-nonpayer-child-account"), body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TransferFunds переводит денежные средства на другой аккаунт (сумма >= 1).
//
// POST v1/account/{accountID}/funds-transfer
func (s *ChildAccountService) TransferFunds(ctx context.Context, accountIDTo, amount int64) error {
	if amount < 1 {
		return validationErrorf("сумма перевода должна быть не менее 1, передано %d", amount)
	}
	body := map[string]any{"accountIdTo": accountIDTo, "amount": amount}
	_, err := s.tr.do(ctx, "POST", s.accountPath("funds-transfer"), body, nil)
	return err
}

// TransferBonus переводит бонусы на другой аккаунт (сумма >= 1).
//
// POST v1/account/{accountID}/bonus-transfer
func (s *ChildAccountService) TransferBonus(ctx context.Context, accountIDTo, amount int64) error {
	if amount < 1 {
		return validationErrorf("сумма перевода должна быть не менее 1, передано %d", amount)
	}
	body := map[string]any{"accountIdTo": accountIDTo, "amount": amount}
	_, err := s.tr.do(ctx, "POST", s.accountPath("bonus-transfer"), body, nil)
	return err
}
