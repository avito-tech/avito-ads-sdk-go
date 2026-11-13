package avito

import "strings"

// Версия SDK совпадает с версией публичного API (info.version в swagger).
const (
	// APIVersion — версия публичного API Авито Реклама (swagger), с которой
	// синхронизирован SDK.
	APIVersion = "1.4.0"
	// Version — версия SDK. Единственный источник версии в коде: из неё
	// собирается метка User-Agent.
	Version = "1.4.0"
)

// userAgentProduct — метка SDK в заголовке User-Agent. По ней сервер
// определяет, что запрос отправлен из Go SDK.
const userAgentProduct = "avito-ads-sdk-go/"

// buildUserAgent возвращает значение заголовка User-Agent.
//
// Без пользовательского значения — метка SDK "avito-ads-sdk-go/<Version>".
// Пользовательское значение не заменяет метку: она дописывается через пробел
// ("my-app/1.0 avito-ads-sdk-go/1.4.0"). Если метка уже есть в пользовательском
// значении (без учёта регистра), второй раз она не добавляется.
func buildUserAgent(custom string) string {
	sdk := userAgentProduct + Version
	custom = strings.TrimSpace(custom)
	if custom == "" {
		return sdk
	}
	if strings.Contains(strings.ToLower(custom), userAgentProduct) {
		return custom
	}
	return custom + " " + sdk
}
