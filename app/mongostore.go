package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/go-multierror"
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

func (s *MongoStore[T]) Create(ctx context.Context, val *T) *Error {
	if errs := s.rdef.Validate(val); len(errs) > 0 {
		var merr *multierror.Error

		for _, err := range errs {
			merr = multierror.Append(merr, err)
		}

		return NewInvalidResourceError(s.rdef.Name, merr)
	}

	if _, err := s.col.InsertOne(ctx, val); err != nil {
		action := fmt.Sprintf("insert %s", s.rdef.Name)

		return NewActionFailedError(action, err)
	}

	return nil
}

func (s *MongoStore[T]) CreateFromJSON(ctx context.Context, r io.Reader) *Error {
	var val T

	if err := json.NewDecoder(r).Decode(&val); err != nil {
		return NewActionFailedError("decode JSON", err)
	}

	return s.Create(ctx, &val)
}

func (s *MongoStore[T]) Delete(ctx context.Context, key string) *Error {
	result, err := s.col.DeleteOne(ctx, bson.D{{"_id", key}})
	if err != nil {
		action := fmt.Sprintf("delete %s with %s '%s'", s.rdef.Name, s.rdef.KeyName, key)

		return NewActionFailedError(action, err)
	}

	if result.DeletedCount == 0 {
		return s.rdef.NewNotFoundError(key)
	}

	return nil
}

func (s *MongoStore[T]) Get(ctx context.Context, key string) (*T, *Error) {
	var val T

	err := s.col.FindOne(ctx, bson.D{{"_id", key}}).Decode(&val)
	if err == mongo.ErrNoDocuments {
		return nil, s.rdef.NewNotFoundError(key)
	} else if err != nil {
		action := fmt.Sprintf("find %s with %s '%s'", s.rdef.Name, s.rdef.KeyName, key)

		return nil, NewActionFailedError(action, err)
	}

	return &val, nil
}

func (s *MongoStore[T]) GetAll(ctx context.Context) ([]*T, *Error) {
	cursor, err := s.col.Find(ctx, bson.D{})
	if err != nil {
		action := fmt.Sprintf("find all %s", s.rdef.NamePlural)

		return []*T{}, NewActionFailedError(action, err)
	}

	var all []*T

	err = cursor.All(ctx, &all)
	if err != nil {
		action := fmt.Sprintf("decode find results as %s", s.rdef.NamePlural)

		return []*T{}, NewActionFailedError(action, err)
	}

	// for no results
	if all == nil {
		all = []*T{}
	}

	return all, nil
}

func (s *MongoStore[T]) Update(ctx context.Context, val *T) *Error {
	if errs := s.rdef.Validate(val); len(errs) > 0 {
		var merr *multierror.Error

		for _, err := range errs {
			merr = multierror.Append(merr, err)
		}

		return NewInvalidResourceError(s.rdef.Name, merr)
	}

	key := s.rdef.GetKey(val)

	_, err := s.col.ReplaceOne(ctx, bson.D{{"_id", key}}, val)
	if err != nil {
		action := fmt.Sprintf("update %s %s", s.rdef.Name, key)

		return NewActionFailedError(action, err)
	}

	return nil
}

func (s *MongoStore[T]) UpdateFromJSON(ctx context.Context, key string, r io.Reader) *Error {
	var val T

	if err := json.NewDecoder(r).Decode(&val); err != nil {
		return NewActionFailedError("decode JSON", err)
	}

	return s.Update(ctx, &val)
}
