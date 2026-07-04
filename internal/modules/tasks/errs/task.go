package errs

import (
	errs2 "github.com/gopinniped/gotion/internal/shared/errs"
)

var (
	ErrTaskNotFound  = errs2.New(errs2.TypeNotFound, "task not found")
	ErrTaskForbidden = errs2.New(errs2.TypePermission, "you don't have access to this task")
)
