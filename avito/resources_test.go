package avito

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

const acct = "/ads/v1/account/123456789"

func TestCampaignsListSendsEmptyFilterAsObject(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("POST", acct+"/campaigns", jsonHandler(200, map[string]any{
		"total": 2,
		"campaigns": []map[string]any{
			{"id": 1, "name": "Кампания 1", "status": "active"},
			{"id": 2, "name": "Кампания 2", "status": "paused"},
		},
	}, nil))

	c := newTestClient(t, ts)
	page, err := c.Campaigns.List(context.Background(), ListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 2 || len(page.Items) != 2 || page.Items[0].Name != "Кампания 1" {
		t.Errorf("page = %+v", page)
	}
	last := ts.requests[len(ts.requests)-1]
	if last.Method != "POST" || last.Path != acct+"/campaigns" {
		t.Errorf("method/path = %s %s", last.Method, last.Path)
	}
	// Пустой фильтр должен сериализоваться как объект {}.
	if !strings.Contains(last.Body, `"filter":{}`) {
		t.Errorf("body должен содержать \"filter\":{} , получено %s", last.Body)
	}
}

func TestCampaignsIterateWalksPages(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.enqueue("POST", acct+"/campaigns", jsonHandler(200, map[string]any{
		"total": 3, "campaigns": []map[string]any{{"id": 1, "name": "a"}, {"id": 2, "name": "b"}},
	}, nil))
	ts.enqueue("POST", acct+"/campaigns", jsonHandler(200, map[string]any{
		"total": 3, "campaigns": []map[string]any{{"id": 3, "name": "c"}},
	}, nil))

	c := newTestClient(t, ts)
	var names []string
	err := c.Campaigns.Iterate(context.Background(), ListRequest{Limit: 2}, func(cm Campaign) bool {
		names = append(names, cm.Name)
		return true
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(names, ",") != "a,b,c" {
		t.Errorf("names = %v", names)
	}
}

func TestChangeBudgetRejectsZero(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	c := newTestClient(t, ts)
	if err := c.Groups.ChangeBudget(context.Background(), 10, 0); err == nil {
		t.Error("ожидалась ошибка валидации")
	}
}

func TestTransferFundsRejectsBelowOne(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	c := newTestClient(t, ts)
	if err := c.ChildAccount.TransferFunds(context.Background(), 999, 0); err == nil {
		t.Error("ожидалась ошибка валидации")
	}
}

func TestTransferFundsSendsBody(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("POST", acct+"/funds-transfer", jsonHandler(200, map[string]any{}, nil))
	c := newTestClient(t, ts)
	if err := c.ChildAccount.TransferFunds(context.Background(), 987654321, 5000); err != nil {
		t.Fatal(err)
	}
	last := ts.requests[len(ts.requests)-1]
	var body map[string]any
	_ = json.Unmarshal([]byte(last.Body), &body)
	if body["accountIdTo"].(float64) != 987654321 || body["amount"].(float64) != 5000 {
		t.Errorf("body = %v", body)
	}
}

func TestStatisticsRejectsPeriodOver100Days(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	c := newTestClient(t, ts)
	if _, err := c.Statistics.Campaign(context.Background(), 1, "2025-01-01", "2025-06-01"); err == nil {
		t.Error("ожидалась ошибка валидации периода")
	}
}

func TestStatisticsCampaignParsesAndSendsDates(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("POST", acct+"/campaigns/1/stats", jsonHandler(200, map[string]any{
		"campaign": map[string]any{
			"id": 1, "name": "Промо",
			"data":      []map[string]any{{"timestamp": "2025-01-01", "views": 100, "clicks": 5}},
			"totalData": map[string]any{"views": 100, "clicks": 5},
		},
		"groups": []any{}, "creatives": []any{},
	}, nil))

	c := newTestClient(t, ts)
	res, err := c.Statistics.Campaign(context.Background(), 1, "2025-01-01", "2025-01-31")
	if err != nil {
		t.Fatal(err)
	}
	if res.Campaign.Name != "Промо" || res.Campaign.TotalData.Views != 100 {
		t.Errorf("res = %+v", res.Campaign)
	}
	last := ts.requests[len(ts.requests)-1]
	if !strings.Contains(last.Body, `"dateFrom":"2025-01-01"`) {
		t.Errorf("body = %s", last.Body)
	}
}

func TestUsersListParsed(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("GET", acct+"/users", jsonHandler(200, map[string]any{
		"total": 1, "users": []map[string]any{{"id": 42, "role": "admin", "hasLoggedIn": true}},
	}, nil))
	c := newTestClient(t, ts)
	users, err := c.Users.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 || users[0].ID != 42 || users[0].Role != "admin" || !users[0].HasLoggedIn {
		t.Errorf("users = %+v", users)
	}
}

func TestAddUserRejectsInvalidRole(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	c := newTestClient(t, ts)
	if err := c.Users.Add(context.Background(), 42, "superhero"); err == nil {
		t.Error("ожидалась ошибка валидации роли")
	}
}

func TestDeleteUserUsesDeleteMethod(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("DELETE", acct+"/delete-user/42", jsonHandler(200, map[string]any{}, nil))
	c := newTestClient(t, ts)
	if err := c.Users.Delete(context.Background(), 42); err != nil {
		t.Fatal(err)
	}
	last := ts.requests[len(ts.requests)-1]
	if last.Method != "DELETE" || last.Path != acct+"/delete-user/42" {
		t.Errorf("method/path = %s %s", last.Method, last.Path)
	}
}

func TestCreateContractViaBuilder(t *testing.T) {
	ts := newTestServer()
	defer ts.close()
	ts.on("POST", "/token", tokenHandler("t"))
	ts.on("POST", acct+"/create-contract", jsonHandler(200, map[string]any{"id": 555}, nil))
	c := newTestClient(t, ts)
	builder := ExternalContract("EXT-1").Advertiser(987).
		CounterpartyType(CounterpartyAdvertiserIntermediary).
		Intermediary(map[string]any{"inn": "1"})
	created, err := c.Contracts.Create(context.Background(), builder)
	if err != nil {
		t.Fatal(err)
	}
	if created.ID != 555 {
		t.Errorf("id = %d", created.ID)
	}
	last := ts.requests[len(ts.requests)-1]
	if !strings.Contains(last.Body, `"type":"external"`) || !strings.Contains(last.Body, `"cid":"EXT-1"`) {
		t.Errorf("body = %s", last.Body)
	}
}
