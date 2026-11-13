package avito

// Тесты синхронизации со swagger 1.4.0: новые поля, nullable-бюджет, политика
// неизвестных значений enum, 409, 207, User-Agent и версия.

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestVersionMatchesOpenAPI(t *testing.T) {
	raw, err := os.ReadFile("../openapi.json")
	if err != nil {
		t.Fatalf("openapi.json: %v", err)
	}
	var spec struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
	}
	if err := json.Unmarshal(raw, &spec); err != nil {
		t.Fatalf("openapi.json: %v", err)
	}
	if spec.Info.Version != APIVersion || Version != APIVersion {
		t.Errorf("версии расходятся: openapi.json=%q APIVersion=%q Version=%q",
			spec.Info.Version, APIVersion, Version)
	}
}

func TestBuildUserAgent(t *testing.T) {
	sdk := "avito-ads-sdk-go/" + Version
	cases := []struct{ custom, want string }{
		{"", sdk},
		{"  ", sdk},
		{"my-app/1.0", "my-app/1.0 " + sdk},
		{"my-app/1.0 avito-ads-sdk-go/1.4.0", "my-app/1.0 avito-ads-sdk-go/1.4.0"},
		{"My-App AVITO-ADS-SDK-GO/9.9.9", "My-App AVITO-ADS-SDK-GO/9.9.9"},
	}
	for _, tc := range cases {
		if got := buildUserAgent(tc.custom); got != tc.want {
			t.Errorf("buildUserAgent(%q) = %q, want %q", tc.custom, got, tc.want)
		}
	}
}

func TestUserAgentSentInTokenAndAPIRequests(t *testing.T) {
	cases := []struct {
		name string
		opts []Option
		want string
	}{
		{"default", nil, "avito-ads-sdk-go/" + Version},
		{"custom", []Option{WithUserAgent("my-app/1.0")}, "my-app/1.0 avito-ads-sdk-go/" + Version},
		{"no-duplicate", []Option{WithUserAgent("my-app/1.0 avito-ads-sdk-go/" + Version)}, "my-app/1.0 avito-ads-sdk-go/" + Version},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := newTestServer()
			defer ts.close()
			ts.on("POST", "/token", tokenHandler("t"))
			ts.on("GET", balancePath, jsonHandler(200, map[string]any{"balance": 1}, nil))
			opts := append([]Option{
				WithHTTPClient(&http.Client{Transport: &rewriteTransport{target: ts.srv.URL}}),
				WithRetryBaseWait(0),
			}, tc.opts...)
			c, err := NewClient("cid", "secret", 123456789, opts...)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := c.Account.GetBalance(context.Background()); err != nil {
				t.Fatal(err)
			}
			if len(ts.requests) != 2 {
				t.Fatalf("ожидалось 2 запроса (токен и баланс), получено %d", len(ts.requests))
			}
			for _, r := range ts.requests {
				if r.UA != tc.want {
					t.Errorf("%s %s: User-Agent = %q, want %q", r.Method, r.Path, r.UA, tc.want)
				}
			}
		})
	}
}

func TestBalanceKopeks(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("GET", balancePath, jsonHandler(200, map[string]any{
		"balance": 150, "balanceKopeks": 15050, "bonusBalance": 5, "bonusBalanceKopeks": 501,
	}, nil))
	c := newTestClient(t, ts)
	b, err := c.Account.GetBalance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if b.BalanceKopeks != 15050 || b.BonusBalanceKopeks != 501 {
		t.Errorf("balance = %+v", b)
	}
}

func TestCampaignBudgetNullable(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("POST", acct+"/campaigns", jsonHandler(200, `{"total":3,"campaigns":[
		{"id":1,"budget":1000,"budgetKopeks":100050,"budgetFrequency":"daily","statusReasons":[],"topStatusReason":""},
		{"id":2,"budgetFrequency":"mixed","statusReasons":[],"topStatusReason":""},
		{"id":3,"budget":0,"budgetKopeks":null,"budgetFrequency":"entirePeriod","performanceAdditionalAgreementID":77,"statusReasons":[],"topStatusReason":""}
	]}`, nil))
	c := newTestClient(t, ts)
	page, err := c.Campaigns.List(context.Background(), ListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 3 {
		t.Fatalf("items = %d", len(page.Items))
	}
	daily, mixed, empty := page.Items[0], page.Items[1], page.Items[2]
	if daily.Budget == nil || *daily.Budget != 1000 || daily.BudgetKopeks == nil || *daily.BudgetKopeks != 100050 {
		t.Errorf("daily: budget=%v budgetKopeks=%v", daily.Budget, daily.BudgetKopeks)
	}
	if daily.BudgetFrequency != CampaignBudgetFrequencyDaily {
		t.Errorf("daily: budgetFrequency=%q", daily.BudgetFrequency)
	}
	if mixed.Budget != nil || mixed.BudgetKopeks != nil || mixed.BudgetFrequency != CampaignBudgetFrequencyMixed {
		t.Errorf("mixed: отсутствующий бюджет должен быть nil: %+v", mixed)
	}
	if empty.Budget == nil || *empty.Budget != 0 || empty.BudgetKopeks != nil {
		t.Errorf("empty: budget=%v budgetKopeks=%v", empty.Budget, empty.BudgetKopeks)
	}
	if mixed.PerformanceAdditionalAgreementID != nil ||
		empty.PerformanceAdditionalAgreementID == nil || *empty.PerformanceAdditionalAgreementID != 77 {
		t.Errorf("performanceAdditionalAgreementID разобран неверно")
	}
}

// Политика неизвестных значений: значения, которых нет в SDK, сохраняются как есть.
func TestUnknownEnumValuesAreTolerated(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("POST", acct+"/campaigns", jsonHandler(200, `{"total":1,"campaigns":[
		{"id":1,"status":"some_future_status","budgetFrequency":"weekly_new","statusReasons":["manual","brand_new_reason"],"topStatusReason":"brand_new_reason"}
	]}`, nil))
	ts.on("POST", acct+"/groups", jsonHandler(200, `{"total":1,"groups":[
		{"id":2,"status":"some_future_status","budgetFrequency":"hourly","format":"newFormat","pace":"turbo",
		 "frequencyRules":[{"count":3,"period":"quarter"}],"statusReasons":["brand_new_reason"],"topStatusReason":""}
	]}`, nil))
	ts.on("POST", acct+"/creatives", jsonHandler(200, `{"total":1,"creatives":[
		{"id":3,"status":"some_future_status","statusReasons":["brand_new_reason"],"topStatusReason":"brand_new_reason"}
	]}`, nil))
	c := newTestClient(t, ts)
	ctx := context.Background()

	cp, err := c.Campaigns.List(ctx, ListRequest{})
	if err != nil {
		t.Fatalf("campaigns: %v", err)
	}
	cm := cp.Items[0]
	if cm.Status != "some_future_status" || cm.BudgetFrequency != "weekly_new" ||
		strings.Join(cm.StatusReasons, ",") != "manual,brand_new_reason" || cm.TopStatusReason != "brand_new_reason" {
		t.Errorf("campaign = %+v", cm)
	}

	gp, err := c.Groups.List(ctx, ListRequest{})
	if err != nil {
		t.Fatalf("groups: %v", err)
	}
	g := gp.Items[0]
	if g.Status != "some_future_status" || g.BudgetFrequency != "hourly" || g.Format != "newFormat" ||
		g.Pace != "turbo" || len(g.FrequencyRules) != 1 || g.FrequencyRules[0].Period != "quarter" ||
		len(g.StatusReasons) != 1 || g.TopStatusReason != "" {
		t.Errorf("group = %+v", g)
	}

	crp, err := c.Creatives.List(ctx, ListRequest{})
	if err != nil {
		t.Fatalf("creatives: %v", err)
	}
	cr := crp.Items[0]
	if cr.Status != "some_future_status" || cr.TopStatusReason != "brand_new_reason" || len(cr.StatusReasons) != 1 {
		t.Errorf("creative = %+v", cr)
	}
}

func TestGroupNewFields(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("POST", acct+"/groups", jsonHandler(200, `{"total":1,"groups":[{
		"id":5,"status":"will_launch_soon","budget":100,"budgetKopeks":10000,"price":25,"priceKopeks":2550,
		"budgetFrequency":"entirePeriod","frequencyRules":[{"count":3,"period":"day"},{"count":10,"period":"week"}],
		"frequencyCount":3,"frequencyPeriod":"day","startASAP":true,"priceControl":"manual","pace":"even",
		"format":"native","audienceType":"autotargeting","scheduleType":"custom",
		"activitySchedule":{"monday":[9,10,11],"sunday":[]},"actions":["pause"],"userID":42,
		"periodStartAt":"2026-11-01T00:00:00Z","periodEndAt":null,"periodTimezone":"Europe/Moscow",
		"targetingLocationIDs":[637640],"targetingInterestIDs":["auto"],
		"statusReasons":["scheduled_start"],"topStatusReason":"scheduled_start"
	}]}`, nil))
	c := newTestClient(t, ts)
	page, err := c.Groups.List(context.Background(), ListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	g := page.Items[0]
	if g.BudgetKopeks != 10000 || g.PriceKopeks != 2550 || g.BudgetFrequency != BudgetFrequencyEntirePeriod {
		t.Errorf("деньги/периодичность: %+v", g)
	}
	if len(g.FrequencyRules) != 2 || g.FrequencyRules[1] != (FrequencyRule{Count: 10, Period: FrequencyPeriodWeek}) {
		t.Errorf("frequencyRules = %+v", g.FrequencyRules)
	}
	if g.FrequencyCount != 3 || g.FrequencyPeriod != FrequencyPeriodDay {
		t.Errorf("frequencyCount/Period = %d/%q", g.FrequencyCount, g.FrequencyPeriod)
	}
	if !g.StartASAP || g.PriceControl != GroupPriceControlManual || g.Pace != PaceEven ||
		g.Format != GroupFormatNative || g.AudienceType != AudienceTypeAutotargeting ||
		g.ScheduleType != ScheduleTypeCustom || g.UserID != 42 || g.PeriodEndAt != "" ||
		g.PeriodTimezone != "Europe/Moscow" || g.Status != GroupStatusWillLaunchSoon {
		t.Errorf("group = %+v", g)
	}
	if g.ActivitySchedule == nil || len(g.ActivitySchedule.Monday) != 3 || g.ActivitySchedule.Tuesday != nil {
		t.Errorf("activitySchedule = %+v", g.ActivitySchedule)
	}
	if len(g.TargetingLocationIDs) != 1 || g.TargetingInterestIDs[0] != "auto" || g.Actions[0] != "pause" {
		t.Errorf("targeting/actions = %+v", g)
	}
}

func TestCreativeAndContractNewFields(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("POST", acct+"/creatives", jsonHandler(200, `{"total":1,"creatives":[{
		"id":7,"status":"unpausing","kktuFields":[{"key":"kktu","value":"30.10.1"}],
		"landingID":null,"originUrl":"https://example.com/?a=1","previewUrl":"https://p","userID":9,
		"statusReasons":[],"topStatusReason":""
	}]}`, nil))
	ts.on("POST", acct+"/contracts", jsonHandler(200, `{"total":1,"contracts":[{
		"id":1,"accountID":2,"type":"service","isAgentActingForPublisher":true,"isCreativeReporter":null
	}]}`, nil))
	c := newTestClient(t, ts)
	ctx := context.Background()
	cp, err := c.Creatives.List(ctx, ListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	cr := cp.Items[0]
	if cr.Status != CreativeStatusUnpausing || len(cr.KKTUFields) != 1 || cr.KKTUFields[0].Value != "30.10.1" ||
		cr.LandingID != nil || cr.OriginURL == "" || cr.PreviewURL != "https://p" || cr.UserID != 9 ||
		cr.StatusReasons == nil {
		t.Errorf("creative = %+v", cr)
	}
	kp, err := c.Contracts.List(ctx, ListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	k := kp.Items[0]
	if k.IsAgentActingForPublisher == nil || !*k.IsAgentActingForPublisher || k.IsCreativeReporter != nil {
		t.Errorf("contract = %+v", k)
	}
}

func TestChangeBudgetAndPriceConflict(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	conflict := jsonHandler(409, map[string]any{"code": "conflict", "message": "группа с автоставкой"}, nil)
	ts.on("POST", acct+"/group/10/change-budget", conflict)
	ts.on("POST", acct+"/group/10/change-price", conflict)
	c := newTestClient(t, ts)
	ctx := context.Background()
	for name, err := range map[string]error{
		"change-budget": c.Groups.ChangeBudget(ctx, 10, 100),
		"change-price":  c.Groups.ChangePrice(ctx, 10, 25),
	} {
		if !IsConflict(err) {
			t.Errorf("%s: ожидался IsConflict, получено %v", name, err)
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) || apiErr.StatusCode != 409 || apiErr.Code != "conflict" {
			t.Errorf("%s: ошибка должна оставаться *APIError: %v", name, err)
		}
		if IsBadRequest(err) || IsNotFound(err) {
			t.Errorf("%s: 409 не должна классифицироваться как 400/404", name)
		}
	}
	if IsConflict(errors.New("x")) || IsConflict(nil) {
		t.Error("IsConflict для не-API ошибки должен быть false")
	}
}

func sandboxInput() SandboxAccountInput {
	return SandboxAccountInput{
		INN: "7712345678", ShortName: "Тест", LongName: "ООО «Тест»", OGRN: "1177746123456",
		LegalAddress: "Москва", LegalType: LegalTypeUL,
		Contact: map[string]any{"name": "Иван", "phone": "+79990000000"},
	}
}

func TestSandboxAccount207WithWarnings(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("POST", acct, jsonHandler(207, map[string]any{
		"accountID": 555, "warnings": []string{"не удалось пополнить баланс"},
	}, nil))
	c := newTestClient(t, ts)
	out, err := c.Account.CreateSandboxAccount(context.Background(), sandboxInput())
	if err != nil {
		t.Fatalf("207 должен считаться успехом: %v", err)
	}
	if out.AccountID != 555 || len(out.Warnings) != 1 {
		t.Errorf("out = %+v", out)
	}
	body := ts.requests[len(ts.requests)-1].Body
	if !strings.Contains(body, `"legalType":"ul"`) || strings.Contains(body, "actualAddress") {
		t.Errorf("тело запроса: legalType обязателен, пустой actualAddress не отправляется: %s", body)
	}
}

func TestSandboxAccountRequiresLegalType(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	c := newTestClient(t, ts)
	in := sandboxInput()
	in.LegalType = ""
	if _, err := c.Account.CreateSandboxAccount(context.Background(), in); !errors.Is(err, ErrValidation) {
		t.Errorf("без legalType ожидалась ErrValidation, получено %v", err)
	}
}

func TestAccountGetSingleRequest(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.enqueue("GET", acct, jsonHandler(200, map[string]any{"account": map[string]any{"inn": "7712345678"}}, nil))
	ts.enqueue("GET", acct, jsonHandler(200, map[string]any{"account": nil}, nil))
	c := newTestClient(t, ts)
	ctx := context.Background()

	acc, err := c.Account.Get(ctx)
	if err != nil || acc.INN != "7712345678" {
		t.Fatalf("acc=%v err=%v", acc, err)
	}
	if _, err := c.Account.Get(ctx); err == nil {
		t.Error("ответ без account должен давать ошибку")
	}
	gets := 0
	for _, r := range ts.requests {
		if r.Method == "GET" && r.Path == acct {
			gets++
		}
	}
	if gets != 2 {
		t.Errorf("Get должен делать один запрос на вызов, всего запросов %d", gets)
	}
}

func TestStatisticsKopeksAndCampaignStatistic(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("POST", acct+"/campaigns/555/stats", jsonHandler(200, `{
		"campaign":{"id":555,"name":"К","paymentModel":"CPC","campaignType":"textImage","data":[],
			"totalData":{"views":10,"clicks":1,"ctr":0.1,"spend":12,"spendBonus":0,"spendKopeks":1234,"spendBonusKopeks":5,"cpm":1,"cpc":1}},
		"groups":[{"id":1,"name":"Г","data":[],"totalData":{"spend":1,"spendKopeks":150}}],
		"creatives":[{"id":2,"groupId":1,"name":"Кр","data":[],"totalData":{}}]
	}`, nil))
	c := newTestClient(t, ts)
	st, err := c.Statistics.Campaign(context.Background(), 555, "2026-11-01", "2026-11-10")
	if err != nil {
		t.Fatal(err)
	}
	if st.Campaign.PaymentModel != PaymentModelCPC || st.Campaign.TotalData.SpendKopeks != 1234 ||
		st.Campaign.TotalData.SpendBonusKopeks != 5 || st.Groups[0].TotalData.SpendKopeks != 150 ||
		st.Creatives[0].GroupID != 1 {
		t.Errorf("stats = %+v", st)
	}
}
