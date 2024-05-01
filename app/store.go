package app

import (
	"context"
)

type Store[T any] interface {
	RDef() *ResourceDef[T]

	Create(ctx context.Context, val *T) Error
	Delete(ctx context.Context, key string) Error
	Get(ctx context.Context, key string) (*T, Error)
	GetAll(ctx context.Context) ([]*T, Error)
	Update(ctx context.Context, val *T) Error
}
