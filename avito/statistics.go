package avito

import (
	"context"
	"fmt"
	"time"
)

const maxStatsPeriodDays = 100

// StatisticsService — методы получения статистики.
//
// Максимальная длина периода — 100 дней; даты в формате YYYY-MM-DD.
type StatisticsService struct{ baseService }

// Campaign возвращает статистику по кампании со сводками по группам и креативам.
//
// POST v1/account/{accountID}/campaigns/{campaignID}/stats
func (s *StatisticsService) Campaign(ctx context.Context, campaignID int64, dateFrom, dateTo string) (*CampaignStatisticResult, error) {
	if err := assertPeriod(dateFrom, dateTo); err != nil {
		return nil, err
	}
	body := map[string]any{"dateFrom": dateFrom, "dateTo": dateTo}
	path := s.accountPath(fmt.Sprintf("campaigns/%d/stats", campaignID))
	var out CampaignStatisticResult
	if _, err := s.tr.do(ctx, "POST", path, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Groups возвращает статистику по группам кампании.
//
// POST v1/account/{accountID}/campaigns/{campaignID}/groups/stats
func (s *StatisticsService) Groups(ctx context.Context, campaignID int64, dateFrom, dateTo string, groupIDs []int64) ([]EntityStatistic, error) {
	if err := assertPeriod(dateFrom, dateTo); err != nil {
		return nil, err
	}
	body := map[string]any{"dateFrom": dateFrom, "dateTo": dateTo, "groupIDs": groupIDs}
	path := s.accountPath(fmt.Sprintf("campaigns/%d/groups/stats", campaignID))
	var out struct {
		Groups []EntityStatistic `json:"groups"`
	}
	if _, err := s.tr.do(ctx, "POST", path, body, &out); err != nil {
		return nil, err
	}
	return out.Groups, nil
}

// Creatives возвращает статистику по креативам кампании.
//
// POST v1/account/{accountID}/campaigns/{campaignID}/creatives/stats
func (s *StatisticsService) Creatives(ctx context.Context, campaignID int64, dateFrom, dateTo string, creativeIDs []int64) ([]EntityStatistic, error) {
	if err := assertPeriod(dateFrom, dateTo); err != nil {
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
