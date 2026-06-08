package avito

import "context"

// ContractService — методы работы с договорами.
type ContractService struct{ baseService }

// Create создаёт договор из тела, собранного ContractBuilder.
//
// POST v1/account/{accountID}/create-contract
func (s *ContractService) Create(ctx context.Context, builder *ContractBuilder) (*CreatedContract, error) {
	body, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return s.CreateRaw(ctx, body)
}

// CreateRaw создаёт договор из готового тела запроса (без клиентской валидации).
func (s *ContractService) CreateRaw(ctx context.Context, body map[string]any) (*CreatedContract, error) {
	var out CreatedContract
	if _, err := s.tr.do(ctx, "POST", s.accountPath("create-contract"), body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List возвращает страницу договоров.
//
// POST v1/account/{accountID}/contracts
func (s *ContractService) List(ctx context.Context, req ListRequest) (PaginatedResult[Contract], error) {
	return listResource[Contract](ctx, s.tr, s.accountPath("contracts"), "contracts", req)
}

// Iterate перебирает все договоры по всем страницам.
func (s *ContractService) Iterate(ctx context.Context, req ListRequest, yield func(Contract) bool) error {
	return iterateAll(ctx, req, func(ctx context.Context, r ListRequest) (PaginatedResult[Contract], error) {
		return s.List(ctx, r)
	}, yield)
}
