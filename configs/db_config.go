package configs

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(ctx context.Context, mconf MongoConfiguration) *mongo.Database {
	const link = ""

	newClient, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf(link, mconf.Username, mconf.Password)))
	if err != nil {
		log.Fatal(err)
	}

	if err := newClient.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	if err := newClient.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	return newClient.Database(mconf.DbName)
}
