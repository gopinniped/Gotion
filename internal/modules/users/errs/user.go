package errs

import (
	errs2 "github.com/gopinniped/gotion/internal/shared/errs"
)

var (
	ErrUserNotFound       = errs2.New(errs2.TypeNotFound, "user not found")
	ErrUserAlreadyExists  = errs2.New(errs2.TypeConflict, "user with this username already exists")
	ErrInvalidPassword    = errs2.New(errs2.TypeValidation, "invalid password")
	ErrInvalidCredentials = errs2.New(errs2.TypeUnauthenticated, "invalid username or password")
	ErrInternal           = errs2.New(errs2.TypeInternal, "internal server error")
)
