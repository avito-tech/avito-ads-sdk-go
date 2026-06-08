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

// Balance — баланс аккаунта (рубли и бонусы, целые числа).
type Balance struct {
	Balance      int64 `json:"balance"`
	BonusBalance int64 `json:"bonusBalance"`
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
}

// Campaign — рекламная кампания.
type Campaign struct {
	ID                    int64  `json:"id"`
	Name                  string `json:"name"`
	Status                string `json:"status"`
	Budget                int64  `json:"budget"`
	PaymentModel          string `json:"paymentModel"`
	CampaignType          string `json:"campaignType"`
	StartDate             string `json:"startDate"`
	EndDate               string `json:"endDate"`
	UserID                int64  `json:"userId"`
	AccountID             int64  `json:"accountId"`
	AdvertiserID          int64  `json:"advertiserId"`
	ContractID            int64  `json:"contractId"`
	ManagerID             int64  `json:"managerID"`
	AdditionalAgreementID int64  `json:"additionalAgreementID"`
	CreatedAt             string `json:"createdAt"`
	UpdatedAt             string `json:"updatedAt"`
}

// Group — группа объявлений.
type Group struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	AccountID    int64  `json:"accountID"`
	CampaignID   int64  `json:"campaignID"`
	Status       string `json:"status"`
	Budget       int64  `json:"budget"`
	Price        int64  `json:"price"`
	PaymentModel string `json:"paymentModel"`
	CampaignType string `json:"campaignType"`
	AdvertiserID int64  `json:"advertiserID"`
	ManagerID    int64  `json:"managerID"`
	HaveCreative bool   `json:"haveCreative"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// Creative — креатив.
type Creative struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	ButtonText   string         `json:"buttonText"`
	Link         string         `json:"link"`
	Status       string         `json:"status"`
	GroupID      int64          `json:"groupID"`
	CampaignID   int64          `json:"campaignID"`
	AccountID    int64          `json:"accountID"`
	AdvertiserID int64          `json:"advertiserID"`
	ManagerID    int64          `json:"managerID"`
	PaymentModel string         `json:"paymentModel"`
	CampaignType string         `json:"campaignType"`
	LegalInfo    map[string]any `json:"legalInfo"`
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
}

// CreatedChildAccount — результат создания дочернего аккаунта.
type CreatedChildAccount struct {
	AccountID    int64  `json:"accountID"`
	ClientKey    string `json:"clientKey"`
	ClientSecret string `json:"clientSecret"`
}

// StatsData — точка статистики (за день или агрегированная).
type StatsData struct {
	Timestamp     string  `json:"timestamp"`
	Views         int64   `json:"views"`
	Clicks        int64   `json:"clicks"`
	CTR           float64 `json:"ctr"`
	Spend         int64   `json:"spend"`
	SpendBonus    int64   `json:"spendBonus"`
	CPM           float64 `json:"cpm"`
	CPC           float64 `json:"cpc"`
	VideoViews25  int64   `json:"videoViews25"`
	VideoViews50  int64   `json:"videoViews50"`
	VideoViews75  int64   `json:"videoViews75"`
	VideoViews100 int64   `json:"videoViews100"`
	Q25           float64 `json:"q25"`
	Q50           float64 `json:"q50"`
	Q75           float64 `json:"q75"`
	VTR           float64 `json:"vtr"`
}

// EntityStatistic — статистика по сущности (кампании, группе или креативу).
type EntityStatistic struct {
	ID           int64       `json:"id"`
	Name         string      `json:"name"`
	GroupID      int64       `json:"groupId"`
	PaymentModel string      `json:"paymentModel"`
	CampaignType string      `json:"campaignType"`
	Data         []StatsData `json:"data"`
	TotalData    StatsData   `json:"totalData"`
}

// CampaignStatisticResult — статистика по кампании со сводками по группам и креативам.
type CampaignStatisticResult struct {
	Campaign  EntityStatistic   `json:"campaign"`
	Groups    []EntityStatistic `json:"groups"`
	Creatives []EntityStatistic `json:"creatives"`
}
