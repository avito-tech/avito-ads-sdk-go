package avito

import (
	"errors"
	"fmt"
)

// Базовые ошибки уровня клиента, проверяемые через errors.Is.
var (
	// ErrConfiguration — некорректная конфигурация клиента.
	ErrConfiguration = errors.New("avito: ошибка конфигурации")
	// ErrValidation — ошибка валидации входных данных на стороне клиента.
	ErrValidation = errors.New("avito: ошибка валидации")
)

// APIError описывает ошибку, возвращённую API Авито Реклама.
//
// Для классификации по HTTP-статусу используйте хелперы IsBadRequest,
// IsAuthentication, IsAccessDenied, IsNotFound, IsRateLimit, IsServerError,
// либо извлекайте APIError через errors.As.
type APIError struct {
	// StatusCode — HTTP-код ответа.
	StatusCode int
	// Code — код ошибки из тела ответа (поле "code"), если есть.
	Code string
	// Message — человекочитаемое сообщение (поле "message").
	Message string
	// RetryAfter — значение заголовка Retry-After в секундах (для 429), если есть.
	RetryAfter int
}

// Error реализует интерфейс error.
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("avito: ошибка API %d (%s): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("avito: ошибка API %d: %s", e.StatusCode, e.Message)
}

// statusError извлекает *APIError из цепочки ошибок и сравнивает статус.
func statusError(err error, status int) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == status
	}
	return false
}

// IsBadRequest сообщает, что ошибка — это HTTP 400.
func IsBadRequest(err error) bool { return statusError(err, 400) }

// IsAuthentication сообщает, что ошибка — это HTTP 401.
func IsAuthentication(err error) bool { return statusError(err, 401) }

// IsAccessDenied сообщает, что ошибка — это HTTP 403.
func IsAccessDenied(err error) bool { return statusError(err, 403) }

// IsNotFound сообщает, что ошибка — это HTTP 404.
func IsNotFound(err error) bool { return statusError(err, 404) }

// IsConflict сообщает, что ошибка — это HTTP 409 (конфликт состояния).
//
// API возвращает 409 в GroupService.ChangeBudget и GroupService.ChangePrice,
// когда изменение нельзя применить в текущем состоянии группы. Ошибка
// остаётся *APIError: обработчики через errors.As продолжают её ловить.
func IsConflict(err error) bool { return statusError(err, 409) }

// IsRateLimit сообщает, что ошибка — это HTTP 429 (превышен лимит запросов).
func IsRateLimit(err error) bool { return statusError(err, 429) }

// IsServerError сообщает, что ошибка — это HTTP 5xx.
func IsServerError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode >= 500
	}
	return false
}

// validationErrorf создаёт ошибку валидации, обёрнутую вокруг ErrValidation.
func validationErrorf(format string, args ...any) error {
	return fmt.Errorf("%w: "+format, append([]any{ErrValidation}, args...)...)
}

// configErrorf создаёт ошибку конфигурации, обёрнутую вокруг ErrConfiguration.
func configErrorf(format string, args ...any) error {
	return fmt.Errorf("%w: "+format, append([]any{ErrConfiguration}, args...)...)
}
