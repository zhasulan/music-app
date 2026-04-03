package errors

type ErrorCode string

const (
	CodeInvalidRequest ErrorCode = "invalid_request"
	CodeUnauthorized   ErrorCode = "unauthorized"
	CodeNotFound       ErrorCode = "not_found"
	CodeInternal       ErrorCode = "internal_error"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func New(code ErrorCode, message string) ErrorResponse {
	return ErrorResponse{Error: ErrorBody{Code: code, Message: message}}
}
