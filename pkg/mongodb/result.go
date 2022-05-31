package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type WriteRes interface {
	mongo.InsertOneResult | mongo.InsertManyResult | mongo.UpdateResult | mongo.DeleteResult | mongo.BulkWriteResult
}

func RetW[W WriteRes](w *W, e error) *W {
	if e != nil {
		panic(e)
	}
	return w
}

func RetA[T interface{}](c *mongo.Cursor, e error) T {
	if e != nil {
		panic(e)
	}
	var t T
	if err := c.All(context.TODO(), &t); err != nil {
		panic(err)
	}
	return t
}

func RetD[T interface{}](s *mongo.SingleResult) (T, bool) {
	var t T
	if err := s.Decode(&t); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return t, false
		} else {
			panic(err)
		}
	}
	return t, true
}
