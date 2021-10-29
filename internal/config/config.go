package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Cfg struct {
	Client *mongo.Client
	Env    Env
}
type Env struct {
	MongoAddr      string `env:"MONGO_ADDR"`
	MongoUser      string `env:"MONGO_USER"`
	MongoPwd       string `env:"MONGO_PWD"`
	MongoDBName    string `env:"MONGO_DB_NAME"`
	MongoDBTimeout string `env:"MONGO_DB_TIMEOUT"`
	CollectionUser string `env:"COLLECTION_USER"`
}

func InitConfig() Cfg {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	env := initEvironment()
	return Cfg{
		Client: initConnection(env),
		Env:    env,
	}
}

func initConnection(env Env) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := fmt.Sprintf("mongodb://%v:%v@%v", env.MongoUser, env.MongoPwd, env.MongoAddr)
	fmt.Printf("db uri: %v\n", uri)
	//connect mongodb
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	//ping mongodb
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	return client
}

func initEvironment() Env {
	var e Env
	if err := env.Parse(&e); err != nil {
		panic(err)
	}
	return e
}

func (cfg Cfg) Free() {
	if cfg.Client != nil {
		err := cfg.Client.Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}
}
