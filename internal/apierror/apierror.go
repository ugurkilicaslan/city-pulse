package apierror

import "net/http"

// APIError — standart HTTP hata yapısı
type APIError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string { return e.Message }

// Hazır hata tipleri
var (
	ErrNotFound     = &APIError{Status: http.StatusNotFound, Code: "NOT_FOUND", Message: "kaynak bulunamadı"}
	ErrUnauthorized = &APIError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "yetkisiz erişim"}
	ErrForbidden    = &APIError{Status: http.StatusForbidden, Code: "FORBIDDEN", Message: "bu işlem için yetkiniz yok"}
	ErrBadRequest   = &APIError{Status: http.StatusBadRequest, Code: "BAD_REQUEST", Message: "geçersiz istek"}
	ErrUpstream     = &APIError{Status: http.StatusBadGateway, Code: "UPSTREAM_ERROR", Message: "harici API yanıt vermedi"}
	ErrInternal     = &APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "sunucu hatası"}
	ErrRateLimit    = &APIError{Status: http.StatusTooManyRequests, Code: "RATE_LIMITED", Message: "çok fazla istek — lütfen yavaşlayın"}
	ErrConflict     = &APIError{Status: http.StatusConflict, Code: "CONFLICT", Message: "kaynak zaten mevcut"}
)

func New(status int, code, message string) *APIError {
	return &APIError{Status: status, Code: code, Message: message}
}

func Wrap(err error, status int, code string) *APIError {
	return &APIError{Status: status, Code: code, Message: err.Error()}
}
