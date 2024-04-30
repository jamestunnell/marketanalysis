package api

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindOne[T any](
	ctx context.Context,
	key string,
	res *Resource[T],
	col *mongo.Collection,
) (*T, *HTTPErr) {
	var val T

	err := col.FindOne(ctx, bson.D{{"_id", key}}).Decode(&val)
	if err == mongo.ErrNoDocuments {
		err = fmt.Errorf("%s with %s '%s' not found", res.Name, res.KeyName, key)

		return nil, &HTTPErr{Error: err, StatusCode: http.StatusNotFound}
	} else if err != nil {
		err = fmt.Errorf("failed to find %s with %s '%s': %w", res.Name, res.KeyName, key, err)

		return nil, &HTTPErr{Error: err, StatusCode: http.StatusInternalServerError}
	}

	return &val, nil
}
