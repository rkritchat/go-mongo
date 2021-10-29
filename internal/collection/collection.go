package collection

import (
	"context"
	"go-mongo/internal/config"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CollectionName = "user"
)

func CollectionMigration(client *mongo.Client, env config.Env) {
	migrate(client, env.MongoDBName, env.CollectionUser, context.TODO())
}

func migrate(client *mongo.Client, DBName, collectionName string, ctx context.Context) {
	_ = client.Database(DBName).CreateCollection(ctx, collectionName)
}

func initContext(env config.Env) (context.Context, context.CancelFunc) {
	timeout, _ := strconv.Atoi(env.MongoDBTimeout)
	if timeout == 0 {
		timeout = 30 //default
	}
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}
