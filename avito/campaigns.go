package avito

import "context"

// CampaignService — методы работы с кампаниями.
type CampaignService struct{ baseService }

// List возвращает страницу кампаний.
//
// POST v1/account/{accountID}/campaigns
func (s *CampaignService) List(ctx context.Context, req ListRequest) (PaginatedResult[Campaign], error) {
	return listResource[Campaign](ctx, s.tr, s.accountPath("campaigns"), "campaigns", req)
}

// Iterate перебирает все кампании по всем страницам.
func (s *CampaignService) Iterate(ctx context.Context, req ListRequest, yield func(Campaign) bool) error {
	return iterateAll(ctx, req, func(ctx context.Context, r ListRequest) (PaginatedResult[Campaign], error) {
		return s.List(ctx, r)
	}, yield)
}
