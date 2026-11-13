# Avito Ads Go SDK

[![CI](https://github.com/avito-tech/avito-ads-sdk-go/actions/workflows/ci.yml/badge.svg)](https://github.com/avito-tech/avito-ads-sdk-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/avito-tech/avito-ads-sdk-go.svg)](https://pkg.go.dev/github.com/avito-tech/avito-ads-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/avito-tech/avito-ads-sdk-go)](https://goreportcard.com/report/github.com/avito-tech/avito-ads-sdk-go)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Go SDK для API **Авито Реклама** (Avito Ads API).

- [Об Авито Реклама](http://avito.ru/ads)
- [Возможности API Авито Реклама](https://ads-help.avito.com/external/api)
- [Техническая документация API](https://developers.avito.ru/api-catalog/ads/documentation)

Покрывает все основные методы API: аккаунт и баланс, дочерние аккаунты и переводы средств, рекламодателей, договоры, кампании, группы объявлений, креативы, статистику и управление пользователями. Построен **только на стандартной библиотеке** — никаких внешних зависимостей, поэтому встраивается в любой проект и не конфликтует с вашим набором библиотек.

## Версия

Версия SDK равна версии публичного API (swagger): SDK `1.4.0` реализует swagger `1.4.0`.

```go
fmt.Println(avito.Version, avito.APIVersion) // 1.4.0 1.4.0
```

Спецификация, по которой собран SDK, лежит в [`openapi.json`](openapi.json).

## Возможности

- **Без зависимостей** — только stdlib (`net/http`, `encoding/json`).
- Идиоматичный Go: `context.Context` во всех методах, функциональные опции, дженерики для пагинации.
- OAuth2 `client_credentials`: автоматическое получение, кэширование и обновление токена.
- Окружения **production** и **sandbox**.
- Автоповторы при `429`/`5xx` (экспоненциальный backoff, учёт `Retry-After`) и обновление токена при `401`.
- Типизированные модели и ошибки; классификация через `errors.As` / хелперы `IsNotFound`, `IsRateLimit`, ….
- Постраничная выборка и обход всех страниц одним вызовом.
- Потокобезопасный клиент: создавайте один раз и переиспользуйте.

## Установка

```bash
go get github.com/avito-tech/avito-ads-sdk-go/avito
```

Требуется Go 1.21+.

## Быстрый старт

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/avito-tech/avito-ads-sdk-go/avito"
)

func main() {
	client, err := avito.NewClient(
		"ВАШ_CLIENT_ID",
		"ВАШ_CLIENT_SECRET",
		123456789,            // accountID — идентификатор рекламного аккаунта
		avito.WithProduction(), // или avito.WithSandbox()
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	balance, err := client.Account.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Баланс: %d ₽, бонусы: %d\n", balance.Balance, balance.BonusBalance)
}
```

`accountID` задаётся один раз и подставляется во все запросы; токен привязан к этому аккаунту.

### Из переменных окружения

```go
client, err := avito.NewFromEnv() // читает AVITO_ADS_CLIENT_ID/SECRET/ACCOUNT_ID и др.
```

Переменные: `AVITO_ADS_CLIENT_ID`, `AVITO_ADS_CLIENT_SECRET`, `AVITO_ADS_ACCOUNT_ID` (обязательные); `AVITO_ADS_ENVIRONMENT` (`production`/`sandbox`), `AVITO_ADS_MAX_RETRIES`, `AVITO_ADS_TIMEOUT_SECONDS`, `AVITO_ADS_TOKEN_LEEWAY_SECONDS` (необязательные).

## Опции клиента

```go
client, _ := avito.NewClient(id, secret, accountID,
	avito.WithSandbox(),
	avito.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
	avito.WithMaxRetries(4),
	avito.WithRetryBaseWait(time.Second),
	avito.WithTokenLeeway(60*time.Second),
	avito.WithTokenStorage(myStorage), // своё хранилище токена (например, Redis)
	avito.WithUserAgent("my-app/1.0"),
)
```

### User-Agent

По умолчанию SDK отправляет `User-Agent: avito-ads-sdk-go/<версия>` (например,
`avito-ads-sdk-go/1.4.0`) — и в запрос токена, и в запросы к API. `WithUserAgent` не заменяет
метку SDK, а добавляет своё значение перед ней:

| Опция | Заголовок |
|---|---|
| не задана | `avito-ads-sdk-go/1.4.0` |
| `WithUserAgent("my-app/1.0")` | `my-app/1.0 avito-ads-sdk-go/1.4.0` |
| `WithUserAgent("my-app/1.0 avito-ads-sdk-go/1.4.0")` | без изменений (метка не дублируется) |

| Окружение | Базовый адрес |
|---|---|
| `WithProduction()` (по умолчанию) | `https://api.avito.ru/ads/` |
| `WithSandbox()` | `https://api.avito.ru/ads-sandbox/` |

## Основные операции

```go
ctx := context.Background()

// Аккаунт
account, _ := client.Account.Get(ctx)
balance, _ := client.Account.GetBalance(ctx)

// Дочерние аккаунты и переводы (сумма >= 1)
children, _ := client.ChildAccount.List(ctx)
_ = client.ChildAccount.TransferFunds(ctx, 987654321, 5000)

// Рекламодатели
created, _ := client.Advertisers.Create(ctx, avito.AdvertiserInput{
	INN: "7712345678", ShortName: "ООО Реклама", LongName: "ООО «Реклама»",
	OGRN: "1177746123456", LegalAddress: "…", ActualAddress: "…",
	LegalRole: avito.LegalRoleAdvertiser, LegalType: avito.LegalTypeUL, KPP: "771701001",
})

// Кампании / группы / креативы
_ = client.Groups.ChangeBudget(ctx, 555, 100000) // только ручная ставка, значение >= 1
_ = client.Groups.ChangePrice(ctx, 555, 25)

// Статистика (период <= 100 дней, даты YYYY-MM-DD)
stats, _ := client.Statistics.Campaign(ctx, 555, "2025-01-01", "2025-01-31")
fmt.Println(stats.Campaign.TotalData.Views, stats.Campaign.TotalData.Clicks)

// Пользователи
_ = client.Users.Add(ctx, 42, avito.UserRoleAdmin)
_ = client.Users.Delete(ctx, 42)
```

## Кампании, группы и креативы

### Статус и причины статуса

У кампании, группы и креатива есть `Status`, `StatusReasons` (все коды причин, `[]` — причин нет)
и `TopStatusReason` (главная причина, `""` — не определена).

```go
_ = client.Campaigns.Iterate(ctx, avito.ListRequest{}, func(c avito.Campaign) bool {
	switch c.Status {
	case avito.CampaignStatusWillLaunchSoon, avito.CampaignStatusPartiallyActive:
		fmt.Println(c.Name, "скоро запустится или работает частично")
	case avito.CampaignStatusPaused:
		if c.TopStatusReason != "" {
			fmt.Println(c.Name, "на паузе, причина:", c.TopStatusReason)
		}
		fmt.Println("все причины:", c.StatusReasons)
	default:
		// Новые статусы и коды причин приходят как есть — предусмотрите ветку по умолчанию.
	}
	return true
})
```

SDK не проверяет значения перечислений в ответах: если API вернёт статус, код причины или
периодичность бюджета, которых нет в константах SDK, значение сохранится в поле без ошибки.

### Бюджет кампании: nil вместо нуля

`Campaign.Budget` и `Campaign.BudgetKopeks` — `*int64`. Если у групп кампании разная
периодичность бюджета (`BudgetFrequency == avito.CampaignBudgetFrequencyMixed`), API не
возвращает сумму и поле равно `nil`. У кампании без бюджета приходит `0`.

```go
switch {
case c.Budget == nil:
	fmt.Println("бюджет не суммируется:", c.BudgetFrequency) // mixed
default:
	fmt.Printf("бюджет %d ₽ (%s)\n", *c.Budget, c.BudgetFrequency)
}
```

Точные суммы в копейках: `Balance.BalanceKopeks`, `Balance.BonusBalanceKopeks`,
`StatsData.SpendKopeks`, `StatsData.SpendBonusKopeks`, `Campaign.BudgetKopeks`,
`Group.BudgetKopeks`, `Group.PriceKopeks`.

### Частотные ограничения группы

```go
for _, r := range group.FrequencyRules {
	if r.Period == avito.FrequencyPeriodNoLimit {
		continue
	}
	fmt.Printf("не более %d показов за %s\n", r.Count, r.Period)
}
```

`Group.FrequencyCount` и `Group.FrequencyPeriod` устарели (соответствуют первому правилу).
Периодичность бюджета группы — `Group.BudgetFrequency` (`daily` или `entirePeriod`).

### Конфликт при изменении бюджета и ставки

`ChangeBudget` и `ChangePrice` возвращают 409, если изменение нельзя применить в текущем
состоянии группы:

```go
if err := client.Groups.ChangePrice(ctx, 555, 25); avito.IsConflict(err) {
	// например, у группы автоматическое управление ставкой
}
```

### Песочница: тестовый аккаунт

`LegalType` обязателен. Ответ 207 — успех с предупреждениями: аккаунт создан, но часть
подготовки не выполнена.

```go
created, err := client.Account.CreateSandboxAccount(ctx, avito.SandboxAccountInput{
	INN: "7712345678", ShortName: "Тест", LongName: "ООО «Тест»", OGRN: "1177746123456",
	LegalAddress: "Москва", LegalType: avito.LegalTypeUL,
	Contact: map[string]any{"name": "Иван Иванов", "phone": "+79990000000"},
})
if err == nil && len(created.Warnings) > 0 {
	log.Println("аккаунт подготовлен не полностью:", created.Warnings)
}
```

## Договоры

`ContractBuilder` проверяет обязательные поля по типу договора ещё до запроса.

```go
builder := avito.IntermediaryContract().
	Advertiser(987654321).
	CounterpartyType(avito.CounterpartyDirectWithAdvertiser).
	Subject(avito.ContractSubjectMediation).
	Object(avito.ContractActionCommercial).
	ReportingRequired(true).
	FundsAllocationToPrincipal(false).
	Date("2025-01-15").
	Number("ДА-2025/01").
	Intermediary(map[string]any{
		"shortName": "ООО Реклама", "inn": "7712345678",
		"ogrn": "1177746123456", "kpp": "771701001",
		"legalAddress": "…", "actualAddress": "…", "legalType": avito.LegalTypeUL,
	})

created, err := client.Contracts.Create(ctx, builder)
```

Также доступны `avito.ServiceContract()` и `avito.ExternalContract("CID")`. Доп. соглашение — через `.ParentID(...)` (без `Intermediary`).

## Пагинация

```go
// Постранично
page, _ := client.Campaigns.List(ctx, avito.ListRequest{Limit: 50, Page: 1})
fmt.Println(page.Total, page.APIPointBalance)
for _, c := range page.Items {
	fmt.Println(c.ID, c.Name)
}

// Все страницы одним вызовом (верните false, чтобы остановиться досрочно)
_ = client.Campaigns.Iterate(ctx, avito.ListRequest{}, func(c avito.Campaign) bool {
	fmt.Println(c.Name)
	return true
})
```

### Фильтры

```go
f := avito.NewCampaignsFilter().
	Statuses([]string{avito.CampaignStatusActive, avito.CampaignStatusPartiallyActive}).
	ContractIDs([]int64{10, 20}).
	CreatedAt(avito.DateRange{From: "2025-01-01", To: "2025-01-31"})

page, _ := client.Campaigns.List(ctx, avito.ListRequest{Filter: f.Map()})
```

Пустой фильтр корректно сериализуется в JSON-объект `{}`, как требует API.

## Обработка ошибок

```go
balance, err := client.Account.GetBalance(ctx)
switch {
case err == nil:
	// ок
case avito.IsRateLimit(err):
	var apiErr *avito.APIError
	errors.As(err, &apiErr)
	time.Sleep(time.Duration(apiErr.RetryAfter) * time.Second)
case avito.IsNotFound(err):
	// 404
default:
	var apiErr *avito.APIError
	if errors.As(err, &apiErr) {
		log.Printf("API %d (%s): %s", apiErr.StatusCode, apiErr.Code, apiErr.Message)
	}
}
```

Хелперы: `IsBadRequest`, `IsAuthentication`, `IsAccessDenied`, `IsNotFound`, `IsConflict` (409), `IsRateLimit`, `IsServerError`. Ошибки валидации на стороне клиента оборачивают `avito.ErrValidation` (проверяйте через `errors.Is`).

## Интеграция в популярные фреймворки

Клиент потокобезопасен — создайте его один раз при старте и передавайте в обработчики. Зависимостей у SDK нет, поэтому он не повлияет на ваш `go.mod`.

**net/http** — рабочий пример в [`examples/httpserver`](examples/httpserver).

**Gin**

```go
r := gin.Default()
client, _ := avito.NewFromEnv()
r.GET("/balance", func(c *gin.Context) {
	b, err := client.Account.GetBalance(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": b.Balance})
})
```

**Echo**

```go
e := echo.New()
client, _ := avito.NewFromEnv()
e.GET("/balance", func(c echo.Context) error {
	b, err := client.Account.GetBalance(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]any{"balance": b.Balance})
})
```

**Fiber**

```go
app := fiber.New()
client, _ := avito.NewFromEnv()
app.Get("/balance", func(c *fiber.Ctx) error {
	b, err := client.Account.GetBalance(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}
	return c.JSON(fiber.Map{"balance": b.Balance})
})
```

**Chi**

```go
r := chi.NewRouter()
client, _ := avito.NewFromEnv()
r.Get("/balance", func(w http.ResponseWriter, req *http.Request) {
	b, err := client.Account.GetBalance(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"balance": b.Balance})
})
```

Для DI-контейнеров (`uber/fx`, `google/wire`) предоставьте `*avito.Client` как singleton-провайдер через `avito.NewFromEnv`.

## Примеры

В каталоге [`examples/`](examples): `quickstart`, `pagination`, `contracts`, `httpserver`.

```bash
AVITO_ADS_CLIENT_ID=… AVITO_ADS_CLIENT_SECRET=… AVITO_ADS_ACCOUNT_ID=… \
AVITO_ADS_ENVIRONMENT=sandbox go run ./examples/quickstart
```

## Разработка

```bash
go test ./...           # тесты (используют httptest, без обращений к сети)
go vet ./...
gofmt -l .              # проверка форматирования
go test ./... -race     # тесты с детектором гонок
```

## Лицензия

[MIT](LICENSE).
