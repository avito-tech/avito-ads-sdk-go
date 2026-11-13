package avito

import (
	"context"
	"fmt"
)

// GroupService — методы работы с группами объявлений.
//
// Изменение бюджета и ставки доступно только для групп с ручным управлением ставкой.
type GroupService struct{ baseService }

// List возвращает страницу групп.
//
// POST v1/account/{accountID}/groups
func (s *GroupService) List(ctx context.Context, req ListRequest) (PaginatedResult[Group], error) {
	return listResource[Group](ctx, s.tr, s.accountPath("groups"), "groups", req)
}

// Iterate перебирает все группы по всем страницам.
func (s *GroupService) Iterate(ctx context.Context, req ListRequest, yield func(Group) bool) error {
	return iterateAll(ctx, req, func(ctx context.Context, r ListRequest) (PaginatedResult[Group], error) {
		return s.List(ctx, r)
	}, yield)
}

// ChangeBudget изменяет бюджет группы (значение >= 1, в рублях).
//
// Если изменение нельзя применить в текущем состоянии группы, API отвечает
// 409: проверяйте через IsConflict(err).
//
// POST v1/account/{accountID}/group/{groupID}/change-budget
func (s *GroupService) ChangeBudget(ctx context.Context, groupID, budget int64) error {
	if budget < 1 {
		return validationErrorf("бюджет должен быть не менее 1, передано %d", budget)
	}
	path := s.accountPath(fmt.Sprintf("group/%d/change-budget", groupID))
	_, err := s.tr.do(ctx, "POST", path, map[string]any{"budget": budget}, nil)
	return err
}

// ChangePrice изменяет ставку группы (значение >= 1, в рублях).
//
// Если изменение нельзя применить в текущем состоянии группы, API отвечает
// 409: проверяйте через IsConflict(err).
//
// POST v1/account/{accountID}/group/{groupID}/change-price
func (s *GroupService) ChangePrice(ctx context.Context, groupID, price int64) error {
	if price < 1 {
		return validationErrorf("ставка должна быть не менее 1, передано %d", price)
	}
	path := s.accountPath(fmt.Sprintf("group/%d/change-price", groupID))
	_, err := s.tr.do(ctx, "POST", path, map[string]any{"price": price}, nil)
	return err
}
