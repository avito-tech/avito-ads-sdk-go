package avito

// Account — реквизиты рекламного аккаунта.
type Account struct {
	INN           string         `json:"inn"`
	ShortName     string         `json:"shortName"`
	LongName      string         `json:"longName"`
	OGRN          string         `json:"ogrn"`
	KPP           string         `json:"kpp"`
	LegalAddress  string         `json:"legalAddress"`
	ActualAddress string         `json:"actualAddress"`
	Contact       map[string]any `json:"contact"`
	Manager       map[string]any `json:"manager"`
}

// Balance — баланс аккаунта: целые рубли и точные суммы в копейках.
type Balance struct {
	// Balance — баланс в рублях (целая часть).
	Balance int64 `json:"balance"`
	// BalanceKopeks — баланс в копейках.
	BalanceKopeks int64 `json:"balanceKopeks"`
	// BonusBalance — бонусы в рублях (целая часть).
	BonusBalance int64 `json:"bonusBalance"`
	// BonusBalanceKopeks — бонусы в копейках.
	BonusBalanceKopeks int64 `json:"bonusBalanceKopeks"`
}

// Advertiser — рекламодатель.
type Advertiser struct {
	ID            int64  `json:"id"`
	ShortName     string `json:"shortName"`
	LongName      string `json:"longName"`
	INN           string `json:"inn"`
	OGRN          string `json:"ogrn"`
	KPP           string `json:"kpp"`
	LegalAddress  string `json:"legalAddress"`
	ActualAddress string `json:"actualAddress"`
	LegalType     string `json:"legalType"`
	LegalRole     string `json:"legalRole"`
	AccountID     int64  `json:"accountId"`
}

// Contract — договор.
type Contract struct {
	ID          int64      `json:"id"`
	AccountID   int64      `json:"accountID"`
	Client      Advertiser `json:"client"`
	Contractor  Advertiser `json:"contractor"`
	Type        string     `json:"type"`
	Date        string     `json:"date"`
	Number      string     `json:"number"`
	Subject     string     `json:"subject"`
	Object      string     `json:"object"`
	CID         string     `json:"cid"`
	Description string     `json:"description"`
	ParentID    int64      `json:"parentId"`
	// IsAgentActingForPublisher и IsCreativeReporter — признаки договора
	// (nullable в API: nil, если значение не задано).
	IsAgentActingForPublisher *bool `json:"isAgentActingForPublisher"`
	IsCreativeReporter        *bool `json:"isCreativeReporter"`
}

// Campaign — рекламная кампания.
type Campaign struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	// Status — статус кампании (константы CampaignStatus*). Неизвестное SDK
	// значение сохраняется как есть.
	Status string `json:"status"`
	// StatusReasons — все коды причин текущего статуса; пустой список, если
	// причин нет. Возможны новые коды.
	StatusReasons []string `json:"statusReasons"`
	// TopStatusReason — главная причина статуса; пустая строка, если не
	// определена (даже при непустом StatusReasons).
	TopStatusReason string `json:"topStatusReason"`
	// Budget — сумма бюджетов групп в целых рублях с НДС. nil, если сумма не
	// определена (BudgetFrequency == mixed); 0 — у кампании без бюджета.
	Budget *int64 `json:"budget"`
	// BudgetKopeks — точная сумма бюджетов групп в копейках с НДС; nil — как у Budget.
	BudgetKopeks *int64 `json:"budgetKopeks"`
	// BudgetFrequency — периодичность бюджета: daily, entirePeriod или mixed
	// (группы с разной периодичностью). Неизвестное значение сохраняется как есть.
	BudgetFrequency       CampaignBudgetFrequency `json:"budgetFrequency"`
	PaymentModel          string                  `json:"paymentModel"`
	CampaignType          string                  `json:"campaignType"`
	StartDate             string                  `json:"startDate"`
	EndDate               string                  `json:"endDate"`
	UserID                int64                   `json:"userId"`
	AccountID             int64                   `json:"accountId"`
	AdvertiserID          int64                   `json:"advertiserId"`
	ContractID            int64                   `json:"contractId"`
	ManagerID             int64                   `json:"managerID"`
	AdditionalAgreementID int64                   `json:"additionalAgreementID"`
	// PerformanceAdditionalAgreementID — ID доп. соглашения к исполнительскому
	// договору (nil, если его нет).
	PerformanceAdditionalAgreementID *int64 `json:"performanceAdditionalAgreementID"`
	CreatedAt                        string `json:"createdAt"`
	UpdatedAt                        string `json:"updatedAt"`
}

// Group — группа объявлений.
type Group struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	AccountID  int64  `json:"accountID"`
	CampaignID int64  `json:"campaignID"`
	UserID     int64  `json:"userID"`
	// Status — статус группы (константы GroupStatus*). Неизвестное SDK
	// значение сохраняется как есть.
	Status string `json:"status"`
	// StatusReasons — все коды причин текущего статуса; пустой список, если
	// причин нет. Возможны новые коды.
	StatusReasons []string `json:"statusReasons"`
	// TopStatusReason — главная причина статуса; пустая строка, если не определена.
	TopStatusReason string `json:"topStatusReason"`
	// Actions — доступные действия с группой.
	Actions []string `json:"actions"`
	// Budget — бюджет группы в рублях с НДС.
	Budget int64 `json:"budget"`
	// BudgetKopeks — бюджет группы в копейках с НДС.
	BudgetKopeks int64 `json:"budgetKopeks"`
	// BudgetFrequency — периодичность бюджета: daily или entirePeriod.
	BudgetFrequency BudgetFrequency `json:"budgetFrequency"`
	// Price — ставка в рублях с НДС.
	Price int64 `json:"price"`
	// PriceKopeks — ставка в копейках с НДС.
	PriceKopeks  int64             `json:"priceKopeks"`
	PriceControl GroupPriceControl `json:"priceControl"`
	Pace         Pace              `json:"pace"`
	PaymentModel string            `json:"paymentModel"`
	CampaignType string            `json:"campaignType"`
	Format       GroupFormat       `json:"format"`
	AudienceType AudienceType      `json:"audienceType"`
	// FrequencyRules — правила частотного ограничения показов.
	FrequencyRules []FrequencyRule `json:"frequencyRules"`
	// FrequencyCount — частота показов первого правила.
	//
	// Deprecated: используйте FrequencyRules. Поле сохранено в ответе API для
	// совместимости и соответствует первому правилу.
	FrequencyCount int64 `json:"frequencyCount"`
	// FrequencyPeriod — период первого правила.
	//
	// Deprecated: используйте FrequencyRules.
	FrequencyPeriod FrequencyPeriod `json:"frequencyPeriod"`
	// StartASAP — запускать как можно скорее.
	StartASAP bool `json:"startASAP"`
	// PeriodStartAt и PeriodEndAt — границы периода показов (RFC 3339; пусто, если не заданы).
	PeriodStartAt  string `json:"periodStartAt"`
	PeriodEndAt    string `json:"periodEndAt"`
	PeriodTimezone string `json:"periodTimezone"`
	// ScheduleType — тип расписания; при custom часы показов — в ActivitySchedule.
	ScheduleType     ScheduleType      `json:"scheduleType"`
	ActivitySchedule *ActivitySchedule `json:"activitySchedule"`

	TargetingAgeIDs          []int64  `json:"targetingAgeIDs"`
	TargetingAudienceIDs     []string `json:"targetingAudienceIDs"`
	TargetingFamilyStatusIDs []int64  `json:"targetingFamilyStatusIDs"`
	TargetingGenderIDs       []int64  `json:"targetingGenderIDs"`
	TargetingIncomeIDs       []int64  `json:"targetingIncomeIDs"`
	TargetingInterestIDs     []string `json:"targetingInterestIDs"`
	TargetingLocationIDs     []int64  `json:"targetingLocationIDs"`
	TargetingParentStatusIDs []int64  `json:"targetingParentStatusIDs"`

	AdvertiserID int64  `json:"advertiserID"`
	ManagerID    int64  `json:"managerID"`
	HaveCreative bool   `json:"haveCreative"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// FrequencyRule — правило частотного ограничения: не более Count показов за Period.
type FrequencyRule struct {
	Count  int64           `json:"count"`
	Period FrequencyPeriod `json:"period"`
}

// ActivitySchedule — расписание показов по дням недели: часы (0–23), в которые
// группа активна.
type ActivitySchedule struct {
	Monday    []int `json:"monday"`
	Tuesday   []int `json:"tuesday"`
	Wednesday []int `json:"wednesday"`
	Thursday  []int `json:"thursday"`
	Friday    []int `json:"friday"`
	Saturday  []int `json:"saturday"`
	Sunday    []int `json:"sunday"`
}

// Creative — креатив.
type Creative struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ButtonText  string `json:"buttonText"`
	Link        string `json:"link"`
	// OriginURL — оригинальная ссылка для аналитики.
	OriginURL string `json:"originUrl"`
	// PreviewURL — ссылка на предпросмотр креатива.
	PreviewURL string `json:"previewUrl"`
	// LandingID — ID посадочной страницы (nil, если её нет).
	LandingID *int64 `json:"landingID"`
	// Status — статус креатива (константы CreativeStatus*). Неизвестное SDK
	// значение сохраняется как есть.
	Status string `json:"status"`
	// StatusReasons — все коды причин текущего статуса; пустой список, если
	// причин нет. Возможны новые коды.
	StatusReasons []string `json:"statusReasons"`
	// TopStatusReason — главная причина статуса; пустая строка, если не определена.
	TopStatusReason string `json:"topStatusReason"`
	GroupID         int64  `json:"groupID"`
	CampaignID      int64  `json:"campaignID"`
	AccountID       int64  `json:"accountID"`
	AdvertiserID    int64  `json:"advertiserID"`
	ManagerID       int64  `json:"managerID"`
	UserID          int64  `json:"userID"`
	PaymentModel    string `json:"paymentModel"`
	CampaignType    string `json:"campaignType"`
	// KKTUFields — заполненные поля динамической формы ККТУ.
	KKTUFields []KKTUFieldValue `json:"kktuFields"`
	// LegalInfo — юридическая информация креатива.
	//
	// Deprecated: устарело в API с версии 1.2.0, данные ККТУ — в KKTUFields.
	// Поле сохранено для совместимости.
	LegalInfo map[string]any `json:"legalInfo"`
}

// KKTUFieldValue — заполненное поле динамической формы ККТУ.
type KKTUFieldValue struct {
	// Key — внутренний ключ поля.
	Key string `json:"key"`
	// Value — значение поля, приведённое к строке.
	Value string `json:"value"`
}

// User — пользователь аккаунта.
type User struct {
	ID          int64  `json:"id"`
	Role        string `json:"role"`
	HasLoggedIn bool   `json:"hasLoggedIn"`
}

// ChildAccount — дочерний аккаунт (данные аккаунта вложены в поле account).
type ChildAccount struct {
	Account struct {
		ID        int64  `json:"id"`
		ShortName string `json:"shortName"`
	} `json:"account"`
	Contract map[string]any `json:"contract"`
	Balance  *Balance       `json:"balance"`
}

// CreatedContract — результат создания договора.
type CreatedContract struct {
	ID int64 `json:"id"`
}

// CreatedAdvertiser — результат создания рекламодателя.
type CreatedAdvertiser struct {
	ID int64 `json:"id"`
}

// CreatedAccount — результат создания тестового аккаунта (песочница).
type CreatedAccount struct {
	AccountID int64 `json:"accountID"`
	// Warnings — предупреждения о частично невыполненных операциях подготовки
	// песочницы (ответ 207). Пустой список — аккаунт подготовлен полностью.
	Warnings []string `json:"warnings"`
}

// CreatedChildAccount — результат создания дочернего аккаунта.
type CreatedChildAccount struct {
	AccountID    int64  `json:"accountID"`
	ClientKey    string `json:"clientKey"`
	ClientSecret string `json:"clientSecret"`
}

// StatsData — точка статистики (за день или агрегированная).
type StatsData struct {
	Timestamp  string  `json:"timestamp"`
	Views      int64   `json:"views"`
	Clicks     int64   `json:"clicks"`
	CTR        float64 `json:"ctr"`
	Spend      int64   `json:"spend"`
	SpendBonus int64   `json:"spendBonus"`
	// SpendKopeks и SpendBonusKopeks — потрачено денег и бонусов в копейках.
	SpendKopeks      int64   `json:"spendKopeks"`
	SpendBonusKopeks int64   `json:"spendBonusKopeks"`
	CPM              float64 `json:"cpm"`
	CPC              float64 `json:"cpc"`
	VideoViews25     int64   `json:"videoViews25"`
	VideoViews50     int64   `json:"videoViews50"`
	VideoViews75     int64   `json:"videoViews75"`
	VideoViews100    int64   `json:"videoViews100"`
	Q25              float64 `json:"q25"`
	Q50              float64 `json:"q50"`
	Q75              float64 `json:"q75"`
	VTR              float64 `json:"vtr"`
}

// EntityStatistic — статистика по группе или креативу.
type EntityStatistic struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	// GroupID — ID группы; заполняется только в статистике креативов.
	GroupID int64 `json:"groupId"`
	// PaymentModel не возвращается API для групп и креативов.
	//
	// Deprecated: модель оплаты есть только в статистике кампании — CampaignStatistic.PaymentModel.
	PaymentModel string `json:"paymentModel"`
	// CampaignType не возвращается API для групп и креативов.
	//
	// Deprecated: тип кампании есть только в статистике кампании — CampaignStatistic.CampaignType.
	CampaignType string      `json:"campaignType"`
	Data         []StatsData `json:"data"`
	TotalData    StatsData   `json:"totalData"`
}

// CampaignStatistic — статистика по кампании.
type CampaignStatistic struct {
	ID           int64       `json:"id"`
	Name         string      `json:"name"`
	PaymentModel string      `json:"paymentModel"`
	CampaignType string      `json:"campaignType"`
	Data         []StatsData `json:"data"`
	TotalData    StatsData   `json:"totalData"`
}

// CampaignStatisticResult — статистика по кампании со сводками по группам и креативам.
type CampaignStatisticResult struct {
	Campaign  CampaignStatistic `json:"campaign"`
	Groups    []EntityStatistic `json:"groups"`
	Creatives []EntityStatistic `json:"creatives"`
}
