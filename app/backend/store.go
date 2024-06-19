package backend

import (
	"context"
)

type Store[T Resource] interface {
	GetInfo() *ResourceInfo

	Reset(ctx context.Context) Error
	Create(ctx context.Context, val T) Error
	Delete(ctx context.Context, key string) Error
	Get(ctx context.Context, key string) (T, Error)
	GetAll(ctx context.Context) ([]T, Error)
	GetAllKeys(ctx context.Context) ([]string, Error)
	Upsert(ctx context.Context, val T) Error
}
