package errs

var (
	ErrUserNotFound      = New(TypeNotFound, "user not found")
	ErrUserAlreadyExists = New(TypeConflict, "user with this username already exists")
	ErrInvalidPassword   = New(TypeValidation, "invalid password")
	ErrInvalidCredentials = New(TypeUnauthenticated, "invalid username or password")
	ErrInternal          = New(TypeInternal, "internal server error")
)
