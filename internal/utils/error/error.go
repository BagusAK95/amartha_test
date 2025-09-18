package error

type DefaultResponse struct {
	Message string
	Errors  []string
}

type NotFoundError DefaultResponse       // 404 Not Found
type ForbiddenError DefaultResponse      // 403 Forbidden
type BadRequestError DefaultResponse     // 400 Bad Request
type InternalServerError DefaultResponse // 500 Internal Server Error

func NewNotFoundError(message string, errors ...string) error {
	return &NotFoundError{Message: message, Errors: errors}
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewForbiddenError(message string, errors ...string) error {
	return &ForbiddenError{Message: message, Errors: errors}
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func NewBadRequestError(message string, errors ...string) error {
	return &BadRequestError{Message: message, Errors: errors}
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func NewInternalServerError(message string, errors ...string) error {
	return &InternalServerError{Message: message, Errors: errors}
}

func (e *InternalServerError) Error() string {
	return e.Message
}
