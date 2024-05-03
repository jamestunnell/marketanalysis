package app

import (
	"fmt"
	"strings"
)

type Error interface {
	error

	GetType() ErrorType
	GetMessage() string
	GetDetails() []string
}

type ErrorType int

const (
	ActionFailed ErrorType = iota
	InvalidInput
	NotFound
)

type appErr struct {
	Type    ErrorType
	Message string
	Details []string
}

func (t ErrorType) String() string {
	switch t {
	case ActionFailed:
		return "ActionFailed"
	case InvalidInput:
		return "InvalidInput"
	case NotFound:
		return "NotFound"
	}

	return ""
}

func NewErrNotFound(what string) Error {
	return &appErr{
		Type:    NotFound,
		Message: fmt.Sprintf("%s was not found", what),
		Details: nil,
	}
}

func NewErrInvalidInput(name string, reasons ...string) Error {
	return &appErr{
		Type:    InvalidInput,
		Message: fmt.Sprintf("invalid %s", name),
		Details: reasons,
	}
}

func NewErrActionFailed(action string, reason string) Error {
	return &appErr{
		Type:    ActionFailed,
		Message: fmt.Sprintf("failed to %s", action),
		Details: []string{reason},
	}
}

func (err *appErr) Error() string {
	switch len(err.Details) {
	case 0:
		return err.Message
	case 1:
		return fmt.Sprintf("%s: %s", err.Message, err.Details[0])
	}

	var b strings.Builder

	b.WriteString(err.Message)

	for _, detail := range err.Details {
		b.WriteString(":\n\t*")
		b.WriteString(detail)
	}

	b.WriteString("\n")

	return b.String()
}

func (err *appErr) GetType() ErrorType {
	return err.Type
}

func (err *appErr) GetMessage() string {
	return err.Message
}

func (err *appErr) GetDetails() []string {
	return err.Details
}
