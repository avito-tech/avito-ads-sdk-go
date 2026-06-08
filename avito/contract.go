package avito

// ContractBuilder — конструктор тела запроса для создания договора.
//
// Инкапсулирует правила обязательности полей по типу договора и проверяет их
// в Build() до отправки, чтобы вместо ответа 400 вернуть понятную ошибку.
//
// Правила по типу:
//
//	service      — обязательны subject, isReportingRequired, date, number; запрещён cid.
//	intermediary — обязательны subject, object, isReportingRequired,
//	               isFundsAllocationToPrincipal, date, number; запрещён cid.
//	external     — обязателен cid; запрещён parentId; date/number необязательны.
//
// Поле intermediary (данные исполнителя) обязательно, если не задан parentId,
// и недопустимо, если parentId задан (доп. соглашение).
type ContractBuilder struct {
	data map[string]any
	errs []error
}

// NewContractBuilder создаёт пустой конструктор договора.
func NewContractBuilder() *ContractBuilder {
	return &ContractBuilder{data: map[string]any{}}
}

// ServiceContract создаёт конструктор договора на оказание услуг.
func ServiceContract() *ContractBuilder {
	return NewContractBuilder().Type(ContractTypeService)
}

// IntermediaryContract создаёт конструктор посреднического договора.
func IntermediaryContract() *ContractBuilder {
	return NewContractBuilder().Type(ContractTypeIntermediary)
}

// ExternalContract создаёт конструктор внешнего договора.
func ExternalContract(cid string) *ContractBuilder {
	return NewContractBuilder().Type(ContractTypeExternal).CID(cid)
}

// Advertiser задаёт идентификатор рекламодателя.
func (b *ContractBuilder) Advertiser(id int64) *ContractBuilder {
	b.data["advertiserId"] = id
	return b
}

// Type задаёт тип договора.
func (b *ContractBuilder) Type(t string) *ContractBuilder {
	if !contains(validContractTypes, t) {
		b.errs = append(b.errs, validationErrorf("недопустимый type договора: %q", t))
		return b
	}
	b.data["type"] = t
	return b
}

// CounterpartyType задаёт тип контрагента (поле API description).
func (b *ContractBuilder) CounterpartyType(t string) *ContractBuilder {
	if !contains(validCounterpartyTypes, t) {
		b.errs = append(b.errs, validationErrorf("недопустимый тип контрагента: %q", t))
		return b
	}
	b.data["description"] = t
	return b
}

// Subject задаёт предмет договора.
func (b *ContractBuilder) Subject(s string) *ContractBuilder {
	if !contains(validContractSubjects, s) {
		b.errs = append(b.errs, validationErrorf("недопустимый subject: %q", s))
		return b
	}
	b.data["subject"] = s
	return b
}

// Object задаёт действие по договору.
func (b *ContractBuilder) Object(o string) *ContractBuilder {
	if !contains(validContractActions, o) {
		b.errs = append(b.errs, validationErrorf("недопустимый object: %q", o))
		return b
	}
	b.data["object"] = o
	return b
}

// CID задаёт CID внешнего договора.
func (b *ContractBuilder) CID(cid string) *ContractBuilder {
	b.data["cid"] = cid
	return b
}

// Date задаёт дату в формате YYYY-MM-DD.
func (b *ContractBuilder) Date(date string) *ContractBuilder {
	b.data["date"] = date
	return b
}

// Number задаёт номер договора.
func (b *ContractBuilder) Number(number string) *ContractBuilder {
	b.data["number"] = number
	return b
}

// ReportingRequired задаёт обязанность отчитываться по креативам.
func (b *ContractBuilder) ReportingRequired(v bool) *ContractBuilder {
	b.data["isReportingRequired"] = v
	return b
}

// FundsAllocationToPrincipal задаёт направление средств в сторону принципала.
func (b *ContractBuilder) FundsAllocationToPrincipal(v bool) *ContractBuilder {
	b.data["isFundsAllocationToPrincipal"] = v
	return b
}

// ParentID задаёт ID родительского договора (для доп. соглашения).
func (b *ContractBuilder) ParentID(id int64) *ContractBuilder {
	b.data["parentId"] = id
	return b
}

// Intermediary задаёт данные исполнителя (посредника).
func (b *ContractBuilder) Intermediary(data map[string]any) *ContractBuilder {
	b.data["intermediary"] = data
	return b
}

// Build проверяет правила по типу договора и возвращает тело запроса.
func (b *ContractBuilder) Build() (map[string]any, error) {
	if len(b.errs) > 0 {
		return nil, b.errs[0]
	}
	if err := b.require("advertiserId", "не задан advertiserId"); err != nil {
		return nil, err
	}
	if err := b.require("type", "не задан type договора"); err != nil {
		return nil, err
	}
	if err := b.require("description", "не задан тип контрагента (CounterpartyType)"); err != nil {
		return nil, err
	}

	contractType, _ := b.data["type"].(string)
	_, hasParent := b.data["parentId"]

	if contractType == ContractTypeExternal {
		if err := b.require("cid", "для внешнего договора (external) обязателен cid"); err != nil {
			return nil, err
		}
		if hasParent {
			return nil, validationErrorf("для внешнего договора (external) нельзя передавать parentId")
		}
	}

	if contractType == ContractTypeService || contractType == ContractTypeIntermediary {
		for _, k := range []string{"subject", "isReportingRequired", "date", "number"} {
			if err := b.require(k, "для договора типа "+contractType+" обязательны subject, isReportingRequired, date, number"); err != nil {
				return nil, err
			}
		}
		if _, ok := b.data["cid"]; ok {
			return nil, validationErrorf("для договора типа %s нельзя передавать cid", contractType)
		}
	}

	if contractType == ContractTypeIntermediary {
		if err := b.require("object", "для посреднического договора обязателен object"); err != nil {
			return nil, err
		}
		if err := b.require("isFundsAllocationToPrincipal", "для посреднического договора обязателен isFundsAllocationToPrincipal"); err != nil {
			return nil, err
		}
	}

	_, hasIntermediary := b.data["intermediary"]
	if hasParent && hasIntermediary {
		return nil, validationErrorf("для доп. соглашения (задан parentId) нельзя передавать intermediary")
	}
	if !hasParent && !hasIntermediary {
		return nil, validationErrorf("не задан parentId — необходимо передать данные исполнителя (intermediary)")
	}

	out := make(map[string]any, len(b.data))
	for k, v := range b.data {
		out[k] = v
	}
	return out, nil
}

func (b *ContractBuilder) require(key, message string) error {
	if _, ok := b.data[key]; !ok {
		return validationErrorf("%s", message)
	}
	return nil
}
