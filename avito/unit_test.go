package avito

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestConfigBaseURIs(t *testing.T) {
	c, err := NewClient("id", "secret", 123)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.BaseURL(); got != "https://api.avito.ru/ads/" {
		t.Errorf("prod base = %q", got)
	}
	if got := c.cfg.tokenURL(); got != "https://api.avito.ru/token" {
		t.Errorf("token url = %q", got)
	}
	cs, _ := NewClient("id", "secret", 123, WithSandbox())
	if got := cs.BaseURL(); got != "https://api.avito.ru/ads-sandbox/" {
		t.Errorf("sandbox base = %q", got)
	}
}

func TestConfigValidation(t *testing.T) {
	if _, err := NewClient("", "secret", 123); !errors.Is(err, ErrConfiguration) {
		t.Errorf("empty client id should be ErrConfiguration, got %v", err)
	}
	if _, err := NewClient("id", "secret", 0); !errors.Is(err, ErrConfiguration) {
		t.Errorf("zero account id should be ErrConfiguration, got %v", err)
	}
}

func TestTokenExpiry(t *testing.T) {
	base := time.Unix(1000, 0)
	tok := accessToken{Value: "abc", ExpiresAt: base.Add(86400 * time.Second), TokenType: "Bearer"}
	if tok.authorizationHeader() != "Bearer abc" {
		t.Errorf("header = %q", tok.authorizationHeader())
	}
	if tok.isExpired(60*time.Second, base) {
		t.Error("должен быть действителен")
	}
	if !tok.isExpired(60*time.Second, base.Add((86400-30)*time.Second)) {
		t.Error("должен истечь с учётом leeway")
	}
}

func TestListRequestEmptyFilterSerializesAsObject(t *testing.T) {
	r := ListRequest{}
	norm, err := r.normalize()
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(norm.body())
	// Пустой фильтр должен быть {} , а не null или [].
	if !contains([]string{`{"filter":{},"limit":20,"page":1}`}, string(b)) {
		t.Errorf("body = %s", b)
	}
}

func TestListRequestValidation(t *testing.T) {
	cases := []ListRequest{
		{Limit: 101, Page: 1},
		{Limit: -1, Page: 1},
		{Limit: 20, Page: -1},
	}
	for _, c := range cases {
		if _, err := c.normalize(); !errors.Is(err, ErrValidation) {
			t.Errorf("%+v should fail validation, got %v", c, err)
		}
	}
}

func TestPaginatedResultHasNextPage(t *testing.T) {
	if !(PaginatedResult[int]{Total: 3, Page: 1, Limit: 2}).HasNextPage() {
		t.Error("1*2 < 3 → должна быть следующая страница")
	}
	if (PaginatedResult[int]{Total: 3, Page: 2, Limit: 2}).HasNextPage() {
		t.Error("2*2 >= 3 → следующей страницы нет")
	}
}

func TestFiltersRenameKeys(t *testing.T) {
	f := NewCampaignsFilter().
		IDs([]int64{1, 2}).
		Statuses([]string{CampaignStatusActive}).
		ContractIDs([]int64{10, 20}).
		CreatedAt(DateRange{From: "2025-01-01", To: "2025-01-31"})
	m := f.Map()
	if _, ok := m["contractIDs"]; !ok {
		t.Error("ожидался ключ contractIDs")
	}
	cr, ok := m["createdAt"].(map[string]any)
	if !ok || cr["from"] != "2025-01-01" {
		t.Errorf("createdAt = %v", m["createdAt"])
	}
	cf := NewCreativesFilter().GroupIDs([]int64{5}).CampaignIDs([]int64{7}).Map()
	if _, ok := cf["groupIDs"]; !ok {
		t.Error("ожидался ключ groupIDs")
	}
	if _, ok := cf["campaignIDs"]; !ok {
		t.Error("ожидался ключ campaignIDs")
	}
}

func TestContractBuilderValid(t *testing.T) {
	body, err := ServiceContract().
		Advertiser(987).
		CounterpartyType(CounterpartyDirectWithAdvertiser).
		Subject(ContractSubjectDistribution).
		ReportingRequired(true).
		Date("2025-01-01").
		Number("A-1").
		Intermediary(map[string]any{"inn": "7712345678"}).
		Build()
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if body["type"] != ContractTypeService || body["description"] != CounterpartyDirectWithAdvertiser {
		t.Errorf("body = %v", body)
	}
}

func TestContractBuilderRules(t *testing.T) {
	// service без date
	if _, err := ServiceContract().Advertiser(1).CounterpartyType(CounterpartyDirectWithAdvertiser).
		Subject(ContractSubjectDistribution).ReportingRequired(true).Number("A-1").
		Intermediary(map[string]any{"inn": "1"}).Build(); !errors.Is(err, ErrValidation) {
		t.Error("service без date должен падать")
	}
	// external + parentId
	if _, err := ExternalContract("EXT-1").Advertiser(1).CounterpartyType(CounterpartyAdvertiserIntermediary).
		ParentID(5).Build(); !errors.Is(err, ErrValidation) {
		t.Error("external + parentId должен падать")
	}
	// intermediary без object
	if _, err := IntermediaryContract().Advertiser(1).CounterpartyType(CounterpartyDirectWithAdvertiser).
		Subject(ContractSubjectMediation).ReportingRequired(true).Date("2025-01-01").Number("N-1").
		Intermediary(map[string]any{"inn": "1"}).Build(); !errors.Is(err, ErrValidation) {
		t.Error("intermediary без object должен падать")
	}
	// parentId + intermediary вместе
	if _, err := IntermediaryContract().Advertiser(1).CounterpartyType(CounterpartyDirectWithAdvertiser).
		Subject(ContractSubjectMediation).Object(ContractActionCommercial).ReportingRequired(true).
		FundsAllocationToPrincipal(false).Date("2025-01-01").Number("N-1").
		ParentID(7).Intermediary(map[string]any{"inn": "1"}).Build(); !errors.Is(err, ErrValidation) {
		t.Error("parentId + intermediary должен падать")
	}
	// доп. соглашение (parentId, без intermediary) валидно
	if _, err := IntermediaryContract().Advertiser(1).CounterpartyType(CounterpartyDirectWithAdvertiser).
		Subject(ContractSubjectMediation).Object(ContractActionCommercial).ReportingRequired(true).
		FundsAllocationToPrincipal(true).Date("2025-03-01").Number("DS-1").ParentID(100).Build(); err != nil {
		t.Errorf("доп. соглашение должно быть валидно: %v", err)
	}
	// невалидный enum
	if _, err := NewContractBuilder().Subject("bad").Build(); !errors.Is(err, ErrValidation) {
		t.Error("невалидный subject должен падать")
	}
}
