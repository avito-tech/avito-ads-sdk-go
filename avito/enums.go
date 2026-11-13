package avito

// Перечисления значений API.
//
// Политика неизвестных значений: при разборе ответов значения перечислений не
// проверяются. Если API вернёт значение, которого нет среди констант SDK
// (новый статус, новый код причины, новая периодичность бюджета), оно
// сохраняется в поле как есть, без ошибки. Сравнивайте с константами и
// предусматривайте ветку по умолчанию. Допустимость значений проверяется
// только на входе (роли пользователей, поля договоров, legalType).
//
// Статусы и прежние перечисления — нетипизированные строковые константы (их
// можно передавать в фильтры []string). Перечисления, добавленные в 1.4.0, —
// строковые типы с константами.

// LegalType — тип юридического лица.
const (
	LegalTypeUL = "ul"
	LegalTypeIP = "ip"
)

// LegalRole — роль контрагента в системе ОРД.
const (
	LegalRoleAdvertiser  = "rd"
	LegalRoleAgency      = "ra"
	LegalRoleDistributor = "rr"
)

// UserRole — роль пользователя в аккаунте.
const (
	UserRoleAdmin  = "admin"
	UserRoleViewer = "viewer"
)

// ContractType — тип договора.
const (
	ContractTypeService      = "service"
	ContractTypeIntermediary = "intermediary"
	ContractTypeExternal     = "external"
)

// ContractSubject — предмет договора.
const (
	ContractSubjectOrgDistribution = "org-distribution"
	ContractSubjectMediation       = "mediation"
	ContractSubjectDistribution    = "distribution"
	ContractSubjectRepresentation  = "representation"
	ContractSubjectOther           = "other"
)

// ContractAction — действие по договору (поле object).
const (
	ContractActionDistribution = "distribution"
	ContractActionConclude     = "conclude"
	ContractActionCommercial   = "commercial"
	ContractActionOther        = "other"
)

// ContractCounterpartyType — тип контрагента (поле description).
const (
	CounterpartyDirectWithAdvertiser   = "direct_with_advertiser"
	CounterpartyAdvertiserIntermediary = "advertiser_intermediary"
)

// CampaignPaymentModel — модель оплаты кампании.
const (
	PaymentModelCPM = "CPM"
	PaymentModelCPC = "CPC"
)

// CampaignType — тип кампании.
const (
	CampaignTypeTextImage = "textImage"
	CampaignTypeHTML      = "HTML"
	CampaignTypeVideo     = "video"
)

// CampaignStatus — статус кампании.
const (
	CampaignStatusDraft             = "draft"
	CampaignStatusInModeration      = "in_moderation"
	CampaignStatusModerationFailed  = "moderation_failed"
	CampaignStatusPartialModeration = "partial_moderation"
	CampaignStatusActive            = "active"
	CampaignStatusWillLaunchSoon    = "will_launch_soon"
	CampaignStatusPartiallyActive   = "partially_active"
	CampaignStatusWillStopSoon      = "will_stop_soon"
	CampaignStatusPausing           = "pausing"
	CampaignStatusPaused            = "paused"
	CampaignStatusUnpausing         = "unpausing"
	CampaignStatusStopped           = "stopped"
	CampaignStatusFinished          = "finished"
	CampaignStatusArchived          = "archived"
)

// GroupStatus — статус группы объявлений.
const (
	GroupStatusDraft            = "draft"
	GroupStatusInModeration     = "in_moderation"
	GroupStatusModerationFailed = "moderation_failed"
	GroupStatusWillLaunchSoon   = "will_launch_soon"
	GroupStatusActive           = "active"
	GroupStatusWillStopSoon     = "will_stop_soon"
	GroupStatusPausing          = "pausing"
	GroupStatusPaused           = "paused"
	GroupStatusUnpausing        = "unpausing"
	GroupStatusStopped          = "stopped"
	GroupStatusFinished         = "finished"
	GroupStatusArchived         = "archived"
)

// CreativeStatus — статус креатива.
const (
	CreativeStatusDraft              = "draft"
	CreativeStatusReadyForModeration = "ready_for_moderation"
	CreativeStatusInModeration       = "in_moderation"
	CreativeStatusModerationFailed   = "moderation_failed"
	CreativeStatusErirRegistration   = "erir_registration"
	CreativeStatusWillLaunchSoon     = "will_launch_soon"
	CreativeStatusActive             = "active"
	CreativeStatusPausing            = "pausing"
	CreativeStatusPaused             = "paused"
	CreativeStatusUnpausing          = "unpausing"
	CreativeStatusStopped            = "stopped"
	CreativeStatusFinished           = "finished"
	CreativeStatusArchived           = "archived"
)

// CampaignBudgetFrequency — периодичность бюджета кампании (Campaign.BudgetFrequency).
type CampaignBudgetFrequency string

const (
	// CampaignBudgetFrequencyDaily — у всех групп кампании дневной бюджет.
	CampaignBudgetFrequencyDaily CampaignBudgetFrequency = "daily"
	// CampaignBudgetFrequencyEntirePeriod — у всех групп бюджет на весь период.
	CampaignBudgetFrequencyEntirePeriod CampaignBudgetFrequency = "entirePeriod"
	// CampaignBudgetFrequencyMixed — у групп разная периодичность; в этом
	// случае Campaign.Budget и Campaign.BudgetKopeks равны nil.
	CampaignBudgetFrequencyMixed CampaignBudgetFrequency = "mixed"
)

// BudgetFrequency — периодичность бюджета группы (Group.BudgetFrequency).
type BudgetFrequency string

const (
	// BudgetFrequencyDaily — дневной бюджет.
	BudgetFrequencyDaily BudgetFrequency = "daily"
	// BudgetFrequencyEntirePeriod — бюджет на весь период.
	BudgetFrequencyEntirePeriod BudgetFrequency = "entirePeriod"
)

// FrequencyPeriod — период частотного ограничения показов (FrequencyRule.Period).
type FrequencyPeriod string

const (
	FrequencyPeriodNoLimit FrequencyPeriod = "no_limit"
	FrequencyPeriodDay     FrequencyPeriod = "day"
	FrequencyPeriodWeek    FrequencyPeriod = "week"
	FrequencyPeriodMonth   FrequencyPeriod = "month"
)

// GroupFormat — формат группы объявлений.
type GroupFormat string

const (
	GroupFormatMediaRoot          GroupFormat = "mediaRoot"
	GroupFormatMediaSellerProfile GroupFormat = "mediaSellerProfile"
	GroupFormatMediaPremium       GroupFormat = "mediaPremium"
	GroupFormatMedia              GroupFormat = "media"
	GroupFormatNative             GroupFormat = "native"
	GroupFormatHTMLPremium        GroupFormat = "htmlPremium"
	GroupFormatHTML               GroupFormat = "html"
	GroupFormatVideo              GroupFormat = "video"
)

// AudienceType — тип аудитории группы.
type AudienceType string

const (
	AudienceTypeBasic         AudienceType = "basic"
	AudienceTypeReady         AudienceType = "ready"
	AudienceTypeAutotargeting AudienceType = "autotargeting"
)

// Pace — темп открутки бюджета группы (в swagger — Paces).
type Pace string

const (
	PaceASAP Pace = "asap"
	PaceEven Pace = "even"
)

// GroupPriceControl — способ управления ставкой группы.
type GroupPriceControl string

const (
	GroupPriceControlManual GroupPriceControl = "manual"
	GroupPriceControlAuto   GroupPriceControl = "auto"
)

// ScheduleType — тип расписания показов группы.
type ScheduleType string

const (
	ScheduleTypeAlways        ScheduleType = "always"
	ScheduleTypeWeekdays9To20 ScheduleType = "weekdays9To20"
	ScheduleTypeCustom        ScheduleType = "custom"
)

// contains проверяет наличие значения в наборе допустимых.
func contains(allowed []string, value string) bool {
	for _, a := range allowed {
		if a == value {
			return true
		}
	}
	return false
}

var (
	validUserRoles         = []string{UserRoleAdmin, UserRoleViewer}
	validContractTypes     = []string{ContractTypeService, ContractTypeIntermediary, ContractTypeExternal}
	validContractSubjects  = []string{ContractSubjectOrgDistribution, ContractSubjectMediation, ContractSubjectDistribution, ContractSubjectRepresentation, ContractSubjectOther}
	validContractActions   = []string{ContractActionDistribution, ContractActionConclude, ContractActionCommercial, ContractActionOther}
	validCounterpartyTypes = []string{CounterpartyDirectWithAdvertiser, CounterpartyAdvertiserIntermediary}
)
