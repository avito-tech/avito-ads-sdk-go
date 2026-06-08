package avito

// Перечисления значений API. Значения — обычные строки; для проверки
// допустимости используйте функции assert*-валидации в соответствующих местах.

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
	CampaignStatusPaused            = "paused"
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
	CreativeStatusActive             = "active"
	CreativeStatusPaused             = "paused"
	CreativeStatusStopped            = "stopped"
	CreativeStatusFinished           = "finished"
	CreativeStatusArchived           = "archived"
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
