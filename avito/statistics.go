package avito

import (
	"context"
	"fmt"
	"time"
)

const (
	maxStatsPeriodDays = 100
	maxStatsIDs        = 5000
)

// StatisticsOption — необязательный параметр запроса статистики кампании и групп.
type StatisticsOption func(*statisticsOptions)

type statisticsOptions struct {
	includeConversions bool
}

// WithConversions запрашивает конверсии: в StatsData заполняется поле Conversions.
// Доступно для статистики кампании и групп; для креативов API конверсии не возвращает.
func WithConversions() StatisticsOption {
	return func(o *statisticsOptions) { o.includeConversions = true }
}

// StatisticsService — методы получения статистики.
//
// Максимальная длина периода — 100 дней; даты в формате YYYY-MM-DD. В запросах по группам
// и креативам — от 1 до 5000 ID, повторяющиеся ID учитываются в лимите.
type StatisticsService struct{ baseService }

// Campaign возвращает статистику по кампании со сводками по группам и креативам.
// Опция WithConversions добавляет конверсии кампании и групп.
//
// POST v1/account/{accountID}/campaigns/{campaignID}/stats
func (s *StatisticsService) Campaign(ctx context.Context, campaignID int64, dateFrom, dateTo string, opts ...StatisticsOption) (*CampaignStatisticResult, error) {
	if err := assertPeriod(dateFrom, dateTo); err != nil {
		return nil, err
	}
	body := map[string]any{"dateFrom": dateFrom, "dateTo": dateTo}
	applyStatisticsOptions(body, opts)
	path := s.accountPath(fmt.Sprintf("campaigns/%d/stats", campaignID))
	var out CampaignStatisticResult
	if _, err := s.tr.do(ctx, "POST", path, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Groups возвращает статистику по группам кампании (от 1 до 5000 ID).
// Опция WithConversions добавляет конверсии групп.
//
// POST v1/account/{accountID}/campaigns/{campaignID}/groups/stats
func (s *StatisticsService) Groups(ctx context.Context, campaignID int64, dateFrom, dateTo string, groupIDs []int64, opts ...StatisticsOption) ([]EntityStatistic, error) {
	if err := assertPeriod(dateFrom, dateTo); err != nil {
		return nil, err
	}
	if err := assertStatsIDs("groupIDs", groupIDs); err != nil {
		return nil, err
	}
	body := map[string]any{"dateFrom": dateFrom, "dateTo": dateTo, "groupIDs": groupIDs}
	applyStatisticsOptions(body, opts)
	path := s.accountPath(fmt.Sprintf("campaigns/%d/groups/stats", campaignID))
	var out struct {
		Groups []EntityStatistic `json:"groups"`
	}
	if _, err := s.tr.do(ctx, "POST", path, body, &out); err != nil {
		return nil, err
	}
	return out.Groups, nil
}

// Creatives возвращает статистику по креативам кампании (от 1 до 5000 ID).
//
// POST v1/account/{accountID}/campaigns/{campaignID}/creatives/stats
func (s *StatisticsService) Creatives(ctx context.Context, campaignID int64, dateFrom, dateTo string, creativeIDs []int64) ([]EntityStatistic, error) {
	if err := assertPeriod(dateFrom, dateTo); err != nil {
		return nil, err
	}
	if err := assertStatsIDs("creativeIDs", creativeIDs); err != nil {
		return nil, err
	}
	body := map[string]any{"dateFrom": dateFrom, "dateTo": dateTo, "creativeIDs": creativeIDs}
	path := s.accountPath(fmt.Sprintf("campaigns/%d/creatives/stats", campaignID))
	var out struct {
		Creatives []EntityStatistic `json:"creatives"`
	}
	if _, err := s.tr.do(ctx, "POST", path, body, &out); err != nil {
		return nil, err
	}
	return out.Creatives, nil
}

func applyStatisticsOptions(body map[string]any, opts []StatisticsOption) {
	var o statisticsOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}
	if o.includeConversions {
		body["includeConversions"] = true
	}
}

func assertStatsIDs(field string, ids []int64) error {
	if len(ids) < 1 || len(ids) > maxStatsIDs {
		return validationErrorf("%s должен содержать от 1 до %d ID, передано %d", field, maxStatsIDs, len(ids))
	}
	return nil
}

func assertPeriod(dateFrom, dateTo string) error {
	const layout = "2006-01-02"
	from, err := time.Parse(layout, dateFrom)
	if err != nil {
		return validationErrorf("дата dateFrom должна быть в формате YYYY-MM-DD")
	}
	to, err := time.Parse(layout, dateTo)
	if err != nil {
		return validationErrorf("дата dateTo должна быть в формате YYYY-MM-DD")
	}
	if from.After(to) {
		return validationErrorf("dateFrom не может быть позже dateTo")
	}
	days := int(to.Sub(from).Hours()/24) + 1
	if days > maxStatsPeriodDays {
		return validationErrorf("период статистики не может превышать %d дней, запрошено %d", maxStatsPeriodDays, days)
	}
	return nil
}
