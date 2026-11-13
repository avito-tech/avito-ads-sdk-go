# Changelog

Все заметные изменения в этом проекте документируются в этом файле.

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.0.0/),
проект придерживается [семантического версионирования](https://semver.org/lang/ru/).

## [1.4.0] — 2026-11-13

Синхронизация с публичным swagger 1.4.0. С этой версии версия SDK равна версии
публичного API (swagger): `avito.Version == avito.APIVersion == "1.4.0"`.

### Breaking
- `Campaign.Budget` теперь `*int64` (было `int64`), добавлено `Campaign.BudgetKopeks *int64`.
  Для кампаний с `budgetFrequency = mixed` API не возвращает сумму — поле равно `nil`;
  у кампании без бюджета приходит явный `0`. Раньше оба случая давали `0`.
- `SandboxAccountInput.LegalType` обязателен (как в swagger): без него
  `CreateSandboxAccount` возвращает ошибку `ErrValidation`, поле всегда уходит в теле
  запроса. `ActualAddress` стал необязательным (`omitempty`; без него API берёт `LegalAddress`).
- Пользовательский User-Agent теперь дополняется меткой SDK, а не заменяет её:
  `WithUserAgent("my-app/1.0")` → `my-app/1.0 avito-ads-sdk-go/1.4.0`. Если метка
  уже есть в значении, она не дублируется. Тот же заголовок уходит в запрос токена.
- `CampaignStatisticResult.Campaign` теперь типа `CampaignStatistic` (было `EntityStatistic`).
  Доступ к полям (`stats.Campaign.TotalData…`) не меняется; ломается только код, который
  явно присваивает значение переменной типа `EntityStatistic`.
- Изменения swagger, ломающие клиентов API (не код SDK): у периодичности бюджета группы
  удалено значение `weekly`; во входных данных группы `frequencyCount`/`frequencyPeriod`
  заменены на `frequencyRules`. Методов создания и изменения групп в SDK нет; в ответе
  `frequencyCount`/`frequencyPeriod` сохранены.

### Added
- Константы `avito.APIVersion` и `avito.Version`; метка User-Agent собирается из `Version`.
  `openapi.json` (swagger 1.4.0) в корне репозитория, проверка совпадения версий в CI.
- Суммы в копейках: `Balance.BalanceKopeks`, `Balance.BonusBalanceKopeks`,
  `StatsData.SpendKopeks`, `StatsData.SpendBonusKopeks`, `Group.BudgetKopeks`, `Group.PriceKopeks`.
- Кампания: `BudgetFrequency` (тип `CampaignBudgetFrequency`: `daily`, `entirePeriod`, `mixed`),
  `PerformanceAdditionalAgreementID *int64`.
- `StatusReasons []string` и `TopStatusReason string` у `Campaign`, `Group` и `Creative`.
- Статусы кампании: `CampaignStatusWillLaunchSoon`, `CampaignStatusPartiallyActive`,
  `CampaignStatusWillStopSoon`, `CampaignStatusPausing`, `CampaignStatusUnpausing`.
- Статусы креатива: `CreativeStatusWillLaunchSoon`, `CreativeStatusPausing`, `CreativeStatusUnpausing`.
- Группа: `FrequencyRules []FrequencyRule` (`Count`, `Period FrequencyPeriod`),
  `BudgetFrequency` (тип `BudgetFrequency`: `daily`, `entirePeriod`), а также `Actions`,
  `ActivitySchedule`, `AudienceType`, `Format`, `Pace`, `PeriodStartAt`, `PeriodEndAt`,
  `PeriodTimezone`, `PriceControl`, `ScheduleType`, `StartASAP`, `UserID`, `Targeting*IDs` (8 полей).
- Строковые типы перечислений с константами: `CampaignBudgetFrequency`, `BudgetFrequency`,
  `FrequencyPeriod`, `GroupFormat`, `AudienceType`, `Pace` (в swagger — `Paces`),
  `GroupPriceControl`, `ScheduleType`.
- Креатив: `KKTUFields []KKTUFieldValue`, `LandingID *int64`, `OriginURL`, `PreviewURL`, `UserID`.
- Договор: `IsAgentActingForPublisher *bool`, `IsCreativeReporter *bool`.
- `CreatedAccount.Warnings` — предупреждения песочницы; ответ 207 считается успехом.
- `IsConflict(err)` — ответ 409 у `ChangeBudget`/`ChangePrice`. Ошибка остаётся `*APIError`.
- Тип `CampaignStatistic` для статистики кампании (с `PaymentModel` и `CampaignType`).

### Changed
- Явная политика неизвестных значений: значения перечислений в ответах не проверяются и
  сохраняются как есть (новые статусы, коды причин, периодичность бюджета). Проверка
  допустимости — только на входе. Покрыто тестами.
- `AccountService.Get` делает один запрос: ответ всегда в обёртке `{"account": …}`; без неё
  возвращается ошибка (раньше уходил второй запрос).
- Фильтры `Statuses` кампаний и креативов принимают новые статусы (документация).

### Deprecated
- `Creative.LegalInfo` — устарело в API с 1.2.0, используйте `KKTUFields`.
- `Group.FrequencyCount`, `Group.FrequencyPeriod` — используйте `FrequencyRules`.
- `EntityStatistic.PaymentModel`, `EntityStatistic.CampaignType` — API не возвращает их для
  групп и креативов; для кампании используйте `CampaignStatistic`.

## [0.1.0] - Unreleased

### Added
- Первый релиз Go SDK для API Авито Реклама.
- Сервисы: аккаунт и баланс, дочерние аккаунты и переводы средств/бонусов,
  рекламодатели, договоры (с `ContractBuilder`), кампании, группы, креативы,
  статистика, пользователи.
- OAuth2 `client_credentials` с кэшированием и автоматическим обновлением токена.
- Окружения production и sandbox.
- Автоматические повторы при `429`/`5xx` и обновление токена при `401`.
- Типизированные модели и ошибки (`*APIError`, хелперы `Is*`, sentinel-ошибки
  `ErrConfiguration` и `ErrValidation`).
- Постраничная выборка с дженериками и обход всех страниц через `Iterate`.
- Функциональные опции конфигурации; `context.Context` во всех методах.
- Зависимости отсутствуют — только стандартная библиотека.
- Тесты на `net/http/httptest`; примеры в каталоге `examples`.
