package common

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrBadRequest      = errors.New("bad request")
	ErrInternalFailure = errors.New("internal failure")
	ErrNotFound        = errors.New("not found")
)

type AError struct {
	appError error
	svcError *fiber.Error
}

func (e AError) AppError() error {
	return e.appError
}

func (e AError) SvcError() *fiber.Error {
	return e.svcError
}

func NewError(svcErr *fiber.Error, appErr error) error {
	return AError{
		appError: appErr,
		svcError: svcErr,
	}
}

func (e AError) Error() string {
	return errors.Join(e.svcError, e.appError).Error()
}
