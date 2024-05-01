package app

import (
	"context"
	"io"
)

type Store[T any] interface {
	RDef() *ResourceDef[T]

	Create(ctx context.Context, val *T) *Error
	CreateFromJSON(ctx context.Context, r io.Reader) *Error
	Delete(ctx context.Context, key string) *Error
	Get(ctx context.Context, key string) (*T, *Error)
	GetAll(ctx context.Context) ([]*T, *Error)
	Update(ctx context.Context, val *T) *Error
	UpdateFromJSON(ctx context.Context, key string, r io.Reader) *Error
}
