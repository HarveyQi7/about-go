package mongodb

import (
	"about-go/config"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Disconnector func()

func Connect(connCtx context.Context, dsName string, opts ...*options.ClientOptions) *mongo.Client {
	ds := config.Get().Datasources.Get(dsName)
	urlOpts := options.Client().ApplyURI(ds.Url)
	opts = append(opts, urlOpts)
	conn, err := mongo.Connect(connCtx, opts...)
	if err != nil {
		panic(err)
	}
	return conn
}

func Default(dsNames ...string) (*mongo.Client, Disconnector) {
	ds := config.Get().Datasources.Get(dsNames...)
	opts := options.Client().ApplyURI(ds.Url)
	conn, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	disconn := func() {
		if err := conn.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
	return conn, disconn
}
