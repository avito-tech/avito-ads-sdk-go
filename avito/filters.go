package avito

// DateRange — диапазон дат в формате YYYY-MM-DD.
type DateRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// toMap возвращает представление диапазона для тела запроса.
func (d DateRange) toMap() map[string]any {
	return map[string]any{"from": d.From, "to": d.To}
}

// Filter — общий билдер фильтра, накапливающий пары ключ-значение.
//
// Создавайте через NewFilter() и заполняйте методами Set/SetDateRange, либо
// используйте типобезопасные обёртки (CampaignsFilter и т. п.). В методы list
// можно передавать как *Filter, так и обычный map[string]any.
type Filter struct {
	data map[string]any
}

// NewFilter создаёт пустой фильтр.
func NewFilter() *Filter {
	return &Filter{data: map[string]any{}}
}

// Set задаёт произвольное поле фильтра (имя поля — как в API).
func (f *Filter) Set(key string, value any) *Filter {
	f.data[key] = value
	return f
}

// SetDateRange задаёт поле-диапазон дат.
func (f *Filter) SetDateRange(key string, r DateRange) *Filter {
	f.data[key] = r.toMap()
	return f
}

// Map возвращает словарь фильтра с именами полей API.
func (f *Filter) Map() map[string]any {
	out := make(map[string]any, len(f.data))
	for k, v := range f.data {
		out[k] = v
	}
	return out
}

// CampaignsFilter — типобезопасный билдер фильтра кампаний.
type CampaignsFilter struct{ f *Filter }

// NewCampaignsFilter создаёт фильтр кампаний.
func NewCampaignsFilter() *CampaignsFilter { return &CampaignsFilter{f: NewFilter()} }

func (c *CampaignsFilter) IDs(ids []int64) *CampaignsFilter     { c.f.Set("ids", ids); return c }
func (c *CampaignsFilter) Statuses(s []string) *CampaignsFilter { c.f.Set("statuses", s); return c }
func (c *CampaignsFilter) PaymentModels(m []string) *CampaignsFilter {
	c.f.Set("paymentModels", m)
	return c
}
func (c *CampaignsFilter) CampaignTypes(t []string) *CampaignsFilter {
	c.f.Set("campaignTypes", t)
	return c
}
func (c *CampaignsFilter) Managers(ids []int64) *CampaignsFilter { c.f.Set("managers", ids); return c }
func (c *CampaignsFilter) Advertisers(ids []int64) *CampaignsFilter {
	c.f.Set("advertisers", ids)
	return c
}
func (c *CampaignsFilter) ContractIDs(ids []int64) *CampaignsFilter {
	c.f.Set("contractIDs", ids)
	return c
}
func (c *CampaignsFilter) AdditionalAgreementIDs(ids []int64) *CampaignsFilter {
	c.f.Set("additionalAgreementIDs", ids)
	return c
}
func (c *CampaignsFilter) CreatedAt(r DateRange) *CampaignsFilter {
	c.f.SetDateRange("createdAt", r)
	return c
}
func (c *CampaignsFilter) TimeFrame(r DateRange) *CampaignsFilter {
	c.f.SetDateRange("timeFrame", r)
	return c
}

// Map возвращает словарь фильтра.
func (c *CampaignsFilter) Map() map[string]any { return c.f.Map() }

// GroupsFilter — типобезопасный билдер фильтра групп.
type GroupsFilter struct{ f *Filter }

// NewGroupsFilter создаёт фильтр групп.
func NewGroupsFilter() *GroupsFilter { return &GroupsFilter{f: NewFilter()} }

func (g *GroupsFilter) IDs(ids []int64) *GroupsFilter          { g.f.Set("ids", ids); return g }
func (g *GroupsFilter) CampaignIDs(ids []int64) *GroupsFilter  { g.f.Set("campaignIDs", ids); return g }
func (g *GroupsFilter) Statuses(s []string) *GroupsFilter      { g.f.Set("statuses", s); return g }
func (g *GroupsFilter) Paces(p []string) *GroupsFilter         { g.f.Set("paces", p); return g }
func (g *GroupsFilter) Managers(ids []int64) *GroupsFilter     { g.f.Set("managers", ids); return g }
func (g *GroupsFilter) Advertisers(ids []int64) *GroupsFilter  { g.f.Set("advertisers", ids); return g }
func (g *GroupsFilter) PaymentModels(m []string) *GroupsFilter { g.f.Set("paymentModels", m); return g }
func (g *GroupsFilter) TimeFrame(r DateRange) *GroupsFilter {
	g.f.SetDateRange("timeFrame", r)
	return g
}

// Map возвращает словарь фильтра.
func (g *GroupsFilter) Map() map[string]any { return g.f.Map() }

// CreativesFilter — типобезопасный билдер фильтра креативов.
type CreativesFilter struct{ f *Filter }

// NewCreativesFilter создаёт фильтр креативов.
func NewCreativesFilter() *CreativesFilter { return &CreativesFilter{f: NewFilter()} }

func (c *CreativesFilter) IDs(ids []int64) *CreativesFilter      { c.f.Set("ids", ids); return c }
func (c *CreativesFilter) GroupIDs(ids []int64) *CreativesFilter { c.f.Set("groupIDs", ids); return c }
func (c *CreativesFilter) CampaignIDs(ids []int64) *CreativesFilter {
	c.f.Set("campaignIDs", ids)
	return c
}
func (c *CreativesFilter) PaymentModels(m []string) *CreativesFilter {
	c.f.Set("paymentModels", m)
	return c
}
func (c *CreativesFilter) CampaignTypes(t []string) *CreativesFilter {
	c.f.Set("campaignTypes", t)
	return c
}
func (c *CreativesFilter) Statuses(s []string) *CreativesFilter  { c.f.Set("statuses", s); return c }
func (c *CreativesFilter) Managers(ids []int64) *CreativesFilter { c.f.Set("managers", ids); return c }
func (c *CreativesFilter) Advertisers(ids []int64) *CreativesFilter {
	c.f.Set("advertisers", ids)
	return c
}
func (c *CreativesFilter) TimeFrame(r DateRange) *CreativesFilter {
	c.f.SetDateRange("timeFrame", r)
	return c
}

// Map возвращает словарь фильтра.
func (c *CreativesFilter) Map() map[string]any { return c.f.Map() }

// AdvertisersFilter — типобезопасный билдер фильтра рекламодателей.
type AdvertisersFilter struct{ f *Filter }

// NewAdvertisersFilter создаёт фильтр рекламодателей.
func NewAdvertisersFilter() *AdvertisersFilter { return &AdvertisersFilter{f: NewFilter()} }

func (a *AdvertisersFilter) IDs(ids []int64) *AdvertisersFilter    { a.f.Set("ids", ids); return a }
func (a *AdvertisersFilter) INNs(inns []string) *AdvertisersFilter { a.f.Set("inns", inns); return a }
func (a *AdvertisersFilter) Roles(roles []string) *AdvertisersFilter {
	a.f.Set("roles", roles)
	return a
}

// Map возвращает словарь фильтра.
func (a *AdvertisersFilter) Map() map[string]any { return a.f.Map() }

// ContractsFilter — типобезопасный билдер фильтра договоров.
type ContractsFilter struct{ f *Filter }

// NewContractsFilter создаёт фильтр договоров.
func NewContractsFilter() *ContractsFilter { return &ContractsFilter{f: NewFilter()} }

func (c *ContractsFilter) IDs(ids []int64) *ContractsFilter    { c.f.Set("ids", ids); return c }
func (c *ContractsFilter) Numbers(n []string) *ContractsFilter { c.f.Set("numbers", n); return c }
func (c *ContractsFilter) Contractors(ids []int64) *ContractsFilter {
	c.f.Set("contractors", ids)
	return c
}
func (c *ContractsFilter) Clients(ids []int64) *ContractsFilter { c.f.Set("clients", ids); return c }

// Map возвращает словарь фильтра.
func (c *ContractsFilter) Map() map[string]any { return c.f.Map() }
