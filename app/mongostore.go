package app

import (
	"context"
	"fmt"

	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStore[T any] struct {
	rdef *ResourceDef[T]
	col  *mongo.Collection
}

func NewMongoStore[T any](
	rdef *ResourceDef[T],
	col *mongo.Collection,
) Store[T] {
	return &MongoStore[T]{
		rdef: rdef,
		col:  col,
	}
}

func (s *MongoStore[T]) RDef() *ResourceDef[T] {
	return s.rdef
}

func (s *MongoStore[T]) Create(ctx context.Context, val *T) Error {
	if key := s.rdef.GetKey(val); key == "" {
		return NewErrInvalidInput(s.rdef.KeyName, "key is empty")
	}

	if errs := s.rdef.Validate(val); len(errs) > 0 {
		reasons := sliceutils.Map(errs, func(err error) string {
			return err.Error()
		})

		return NewErrInvalidInput(s.rdef.Name, reasons...)
	}

	if _, err := s.col.InsertOne(ctx, val); err != nil {
		action := fmt.Sprintf("insert %s", s.rdef.Name)

		return NewErrActionFailed(action, err.Error())
	}

	return nil
}

func (s *MongoStore[T]) Delete(ctx context.Context, key string) Error {
	result, err := s.col.DeleteOne(ctx, bson.D{{"_id", key}})
	if err != nil {
		action := fmt.Sprintf("delete %s with %s '%s'", s.rdef.Name, s.rdef.KeyName, key)

		return NewErrActionFailed(action, err.Error())
	}

	if result.DeletedCount == 0 {
		return s.rdef.NewErrNotFound(key)
	}

	return nil
}

func (s *MongoStore[T]) Get(ctx context.Context, key string) (*T, Error) {
	var val T

	err := s.col.FindOne(ctx, bson.D{{"_id", key}}).Decode(&val)
	if err == mongo.ErrNoDocuments {
		return nil, s.rdef.NewErrNotFound(key)
	} else if err != nil {
		action := fmt.Sprintf("find %s with %s '%s'", s.rdef.Name, s.rdef.KeyName, key)

		return nil, NewErrActionFailed(action, err.Error())
	}

	return &val, nil
}

func (s *MongoStore[T]) GetAll(ctx context.Context) ([]*T, Error) {
	cursor, err := s.col.Find(ctx, bson.D{})
	if err != nil {
		action := fmt.Sprintf("find all %s", s.rdef.NamePlural)

		return []*T{}, NewErrActionFailed(action, err.Error())
	}

	var all []*T

	err = cursor.All(ctx, &all)
	if err != nil {
		action := fmt.Sprintf("decode find results as %s", s.rdef.NamePlural)

		return []*T{}, NewErrActionFailed(action, err.Error())
	}

	// for no results
	if all == nil {
		all = []*T{}
	}

	return all, nil
}

func (s *MongoStore[T]) Update(ctx context.Context, val *T) Error {
	if errs := s.rdef.Validate(val); len(errs) > 0 {
		reasons := sliceutils.Map(errs, func(err error) string {
			return err.Error()
		})

		return NewErrInvalidInput(s.rdef.Name, reasons...)
	}

	key := s.rdef.GetKey(val)

	_, err := s.col.ReplaceOne(ctx, bson.D{{"_id", key}}, val)
	if err != nil {
		action := fmt.Sprintf("update %s %s", s.rdef.Name, key)

		return NewErrActionFailed(action, err.Error())
	}

	return nil
}
