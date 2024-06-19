package backend

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type MongoStore[T Resource] struct {
	info    *ResourceInfo
	col     *mongo.Collection
	newFunc func() T
}

func NewMongoStore[T Resource](
	info *ResourceInfo,
	col *mongo.Collection,
) Store[T] {
	var val T

	elemType := reflect.TypeOf(val).Elem()

	return &MongoStore[T]{
		info: info,
		col:  col,
		newFunc: func() T {
			return reflect.New(elemType).Interface().(T)
		},
	}
}

func (s *MongoStore[T]) GetInfo() *ResourceInfo {
	return s.info
}

func (s *MongoStore[T]) Reset(ctx context.Context) Error {
	if err := s.col.Drop(ctx); err != nil {
		action := fmt.Sprintf("drop %s collection", s.col.Name())

		return NewErrActionFailed(action, err.Error())
	}

	return nil
}

func (s *MongoStore[T]) Create(ctx context.Context, val T) Error {
	key := val.GetKey()
	if key == "" {
		return NewErrInvalidInput(s.info.KeyName, "key is empty")
	}

	// check for duplicate key
	if err := s.col.FindOne(ctx, bson.M{"_id": key}).Err(); err == nil {
		reason := fmt.Sprintf("%s %s already exists", s.info.KeyName, key)

		return NewErrInvalidInput(s.info.Name, reason)
	}

	if errs := val.Validate(); len(errs) > 0 {
		reasons := sliceutils.Map(errs, func(err error) string {
			return err.Error()
		})

		return NewErrInvalidInput(s.info.Name, reasons...)
	}

	if _, err := s.col.InsertOne(ctx, val); err != nil {
		action := fmt.Sprintf("insert %s", s.info.Name)

		return NewErrActionFailed(action, err.Error())
	}

	return nil
}

func (s *MongoStore[T]) Delete(ctx context.Context, key string) Error {
	result, err := s.col.DeleteOne(ctx, bson.M{"_id": key})
	if err != nil {
		action := fmt.Sprintf("delete %s with %s '%s'", s.info.Name, s.info.KeyName, key)

		return NewErrActionFailed(action, err.Error())
	}

	if result.DeletedCount == 0 {
		return s.info.NewErrNotFound(key)
	}

	return nil
}

func (s *MongoStore[T]) Get(ctx context.Context, key string) (T, Error) {
	val := s.newFunc()

	err := s.col.FindOne(ctx, bson.M{"_id": key}).Decode(val)
	if err == mongo.ErrNoDocuments {
		return val, s.info.NewErrNotFound(key)
	} else if err != nil {
		action := fmt.Sprintf("find %s with %s '%s'", s.info.Name, s.info.KeyName, key)

		return val, NewErrActionFailed(action, err.Error())
	}

	return val, nil
}

func (s *MongoStore[T]) GetAllKeys(ctx context.Context) ([]string, Error) {
	results, err := s.col.Distinct(ctx, "_id", bson.D{})
	if err != nil {
		return []string{}, NewErrActionFailed("get all keys", err.Error())
	}

	keys := sliceutils.Map(results, func(result any) string {
		return result.(string)
	})

	return keys, nil
}

func (s *MongoStore[T]) GetAll(ctx context.Context) ([]T, Error) {
	cursor, err := s.col.Find(ctx, bson.D{})
	if err != nil {
		action := fmt.Sprintf("find all %s", s.info.NamePlural)

		return []T{}, NewErrActionFailed(action, err.Error())
	}

	var all []T

	err = cursor.All(ctx, &all)
	if err != nil {
		action := fmt.Sprintf("decode find results as %s", s.info.NamePlural)

		return []T{}, NewErrActionFailed(action, err.Error())
	}

	// for no results
	if all == nil {
		all = []T{}
	}

	return all, nil
}

func (s *MongoStore[T]) Upsert(ctx context.Context, val T) Error {
	if errs := val.Validate(); len(errs) > 0 {
		reasons := sliceutils.Map(errs, func(err error) string {
			return err.Error()
		})

		return NewErrInvalidInput(s.info.Name, reasons...)
	}

	key := val.GetKey()
	opts := options.Update().SetUpsert(true)

	_, err := s.col.UpdateOne(ctx, bson.M{"_id": key}, val, opts)
	if err != nil {
		action := fmt.Sprintf("update %s %s", s.info.Name, key)

		return NewErrActionFailed(action, err.Error())
	}

	return nil
}
