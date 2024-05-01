package app

import "fmt"

type Error struct {
	Code int
	Err  error
}

const (
	ActionFailed int = iota
	InvalidInput
	NotFound
)

func NewActionFailedError(action string, err error) *Error {
	err = fmt.Errorf("failed to %s: %w", action, err)

	return &Error{
		Code: ActionFailed,
		Err:  err,
	}
}

func NewInvalidResourceError(rname string, validationErr error) *Error {
	err := fmt.Errorf("%s is not valid: %w", rname, validationErr)

	return &Error{
		Code: InvalidInput,
		Err:  err,
	}
}

func NewNotFoundError(descr string) *Error {
	return &Error{
		Code: NotFound,
		Err:  fmt.Errorf("%s not found", descr),
	}
}
