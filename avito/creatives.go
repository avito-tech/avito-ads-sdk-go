package avito

import "context"

// CreativeService — методы работы с креативами.
type CreativeService struct{ baseService }

// List возвращает страницу креативов.
//
// POST v1/account/{accountID}/creatives
func (s *CreativeService) List(ctx context.Context, req ListRequest) (PaginatedResult[Creative], error) {
	return listResource[Creative](ctx, s.tr, s.accountPath("creatives"), "creatives", req)
}

// Iterate перебирает все креативы по всем страницам.
func (s *CreativeService) Iterate(ctx context.Context, req ListRequest, yield func(Creative) bool) error {
	return iterateAll(ctx, req, func(ctx context.Context, r ListRequest) (PaginatedResult[Creative], error) {
		return s.List(ctx, r)
	}, yield)
}
