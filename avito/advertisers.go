package avito

import "context"

// AdvertiserService — методы работы с рекламодателями.
type AdvertiserService struct{ baseService }

// AdvertiserInput — параметры создания рекламодателя.
type AdvertiserInput struct {
	INN           string `json:"inn"`
	ShortName     string `json:"shortName"`
	LongName      string `json:"longName"`
	OGRN          string `json:"ogrn"`
	LegalAddress  string `json:"legalAddress"`
	ActualAddress string `json:"actualAddress"`
	LegalRole     string `json:"legalRole"`
	LegalType     string `json:"legalType"`
	KPP           string `json:"kpp,omitempty"`
}

// Create создаёт рекламодателя.
//
// POST v1/account/{accountID}/create-advertiser
func (s *AdvertiserService) Create(ctx context.Context, in AdvertiserInput) (*CreatedAdvertiser, error) {
	var out CreatedAdvertiser
	if _, err := s.tr.do(ctx, "POST", s.accountPath("create-advertiser"), in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List возвращает страницу рекламодателей.
//
// POST v1/account/{accountID}/advertisers
func (s *AdvertiserService) List(ctx context.Context, req ListRequest) (PaginatedResult[Advertiser], error) {
	return listResource[Advertiser](ctx, s.tr, s.accountPath("advertisers"), "advertisers", req)
}

// Iterate перебирает всех рекламодателей по всем страницам.
func (s *AdvertiserService) Iterate(ctx context.Context, req ListRequest, yield func(Advertiser) bool) error {
	return iterateAll(ctx, req, func(ctx context.Context, r ListRequest) (PaginatedResult[Advertiser], error) {
		return s.List(ctx, r)
	}, yield)
}
