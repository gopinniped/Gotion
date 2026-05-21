package errs

var (
	ErrTaskNotFound  = New(TypeNotFound, "task not found")
	ErrTaskForbidden = New(TypePermission, "you don't have access to this task")
)
