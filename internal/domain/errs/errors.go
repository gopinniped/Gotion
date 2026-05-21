package errs

type ErrorType string

const (
	TypeNotFound        ErrorType = "NOT_FOUND"
	TypeConflict        ErrorType = "CONFLICT"
	TypeValidation      ErrorType = "VALIDATION"
	TypeInternal        ErrorType = "INTERNAL"
	TypeUnauthenticated ErrorType = "UNAUTHENTICATED"
	TypePermission      ErrorType = "PERMISSION"
	TypeTooManyRequests ErrorType = "TOO_MANY_REQ"
)

type AppError struct {
	Type    ErrorType
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(t ErrorType, msg string) error {
	return &AppError{Type: t, Message: msg}
}
