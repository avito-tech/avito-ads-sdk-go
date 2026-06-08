package avito

import (
	"context"
	"encoding/json"
)

// listResource выполняет постраничный POST-запрос и разбирает типизированный результат.
//
// Ответ имеет форму {"total": N, "<itemsKey>": [ ... ]}.
func listResource[T any](
	ctx context.Context,
	tr *transport,
	path, itemsKey string,
	req ListRequest,
) (PaginatedResult[T], error) {
	norm, err := req.normalize()
	if err != nil {
		return PaginatedResult[T]{}, err
	}
	var raw map[string]json.RawMessage
	balance, err := tr.do(ctx, "POST", path, norm.body(), &raw)
	if err != nil {
		return PaginatedResult[T]{}, err
	}
	var total int
	if v, ok := raw["total"]; ok {
		_ = json.Unmarshal(v, &total)
	}
	var items []T
	if v, ok := raw[itemsKey]; ok {
		_ = json.Unmarshal(v, &items)
	}
	return PaginatedResult[T]{
		Total:           total,
		Items:           items,
		Page:            norm.Page,
		Limit:           norm.Limit,
		APIPointBalance: balance,
	}, nil
}

// Лимиты постраничной выборки.
const (
	// MaxLimit — максимальный размер страницы.
	MaxLimit = 100
	// DefaultLimit — размер страницы по умолчанию.
	DefaultLimit = 20
)

// ListRequest — параметры постраничного запроса.
type ListRequest struct {
	// Filter — словарь фильтра. Пустой (но не nil) сериализуется в JSON как {}.
	Filter map[string]any
	// Limit — размер страницы (1..100).
	Limit int
	// Page — номер страницы (>= 1).
	Page int
}

// normalize проверяет и нормализует параметры запроса.
//
// Важно: nil-фильтр заменяется пустым НЕ-nil словарём, иначе encoding/json
// сериализует его как null, а API требует объект {}.
func (r ListRequest) normalize() (ListRequest, error) {
	out := r
	if out.Limit == 0 {
		out.Limit = DefaultLimit
	}
	if out.Page == 0 {
		out.Page = 1
	}
	if out.Limit < 1 || out.Limit > MaxLimit {
		return out, validationErrorf("limit должен быть в диапазоне 1..%d, передано %d", MaxLimit, out.Limit)
	}
	if out.Page < 1 {
		return out, validationErrorf("page должен быть >= 1, передано %d", out.Page)
	}
	if out.Filter == nil {
		out.Filter = map[string]any{}
	}
	return out, nil
}

// body формирует тело запроса.
func (r ListRequest) body() map[string]any {
	filter := r.Filter
	if filter == nil {
		filter = map[string]any{}
	}
	return map[string]any{"filter": filter, "limit": r.Limit, "page": r.Page}
}

// PaginatedResult — результат одной страницы выборки.
type PaginatedResult[T any] struct {
	// Total — общее число элементов по всем страницам.
	Total int `json:"total"`
	// Items — элементы текущей страницы.
	Items []T
	// Page — номер текущей страницы.
	Page int
	// Limit — размер страницы.
	Limit int
	// APIPointBalance — остаток лимита (Api-Point-Balance) или -1, если неизвестен.
	APIPointBalance int
}

// HasNextPage сообщает, есть ли ещё страницы.
func (p PaginatedResult[T]) HasNextPage() bool {
	return p.Page*p.Limit < p.Total
}

// iterateAll перебирает все элементы по всем страницам, вызывая yield для каждого.
// Возврат yield == false останавливает обход. Ошибка загрузки страницы прерывает обход.
func iterateAll[T any](
	ctx context.Context,
	req ListRequest,
	loader func(ctx context.Context, req ListRequest) (PaginatedResult[T], error),
	yield func(T) bool,
) error {
	current, err := req.normalize()
	if err != nil {
		return err
	}
	for {
		page, err := loader(ctx, current)
		if err != nil {
			return err
		}
		if len(page.Items) == 0 {
			return nil
		}
		for _, item := range page.Items {
			if !yield(item) {
				return nil
			}
		}
		if !page.HasNextPage() {
			return nil
		}
		current.Page++
	}
}
