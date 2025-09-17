package error

import "fmt"

type DefaultResponse struct {
	Message string
}

type NotFoundError DefaultResponse   // 404 Not Found
type ForbiddenError DefaultResponse  // 403 Forbidden
type BadRequestError DefaultResponse // 400 Bad Request
type InternalError DefaultResponse   // 500 Internal Server Error

func NewNotFoundError(message string) error {
	return &NotFoundError{Message: message}
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewForbiddenError(message string) error {
	return &ForbiddenError{Message: message}
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) error {
	return &BadRequestError{Message: message}
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func NewInternalError(message string) error {
	return &InternalError{Message: fmt.Sprintf("internal server error: %s", message)}
}

func (e *InternalError) Error() string {
	return e.Message
}
