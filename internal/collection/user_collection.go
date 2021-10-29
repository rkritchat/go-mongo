package collection

import (
	"errors"
	"fmt"
	"go-mongo/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntity struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id,omitempty"` //omitempty for get id from mongo
	Email      string             `json:"email" bson:"email"`
	Credential UserCredential     `json:"credential" bson:"credential"`
	Info       UserInfo           `json:"info" bson:"info"`
}

type UserCredential struct {
	Username     string `json:"username" bson:"username"`
	Password     string `json:"password" bson:"password"`
	HintQuestion string `json:"hint_question" bson:"hint_question"`
	HintAwnser   string `json:"hint_awnser" bson:"hint_awnser"`
}

type UserInfo struct {
	Firstname string `json:"firstname" bson:"first_name"`
	Lastname  string `json:"lastname" bson:"last_name"`
	Age       int    `json:"age" bson:"age"`
	Addreess  string `json:"address" bson:"address"`
}

type User interface {
	DeleteByEmail(email string) error
	Create(entity *UserEntity) error //pass pointer becuase need to update field _id in struct
	UpdateById(entity UserEntity) error
	UpdateByCondition(entity UserEntity) error
	FindOneById(id primitive.ObjectID) (*UserEntity, error)
}

type user struct {
	client         *mongo.Client
	DBName         string
	collectionName string
	env            config.Env
}

func NewUser(client *mongo.Client, env config.Env) User {
	return &user{
		client:         client,
		DBName:         env.MongoDBName,
		collectionName: env.CollectionUser,
	}
}

func (c user) DeleteByEmail(email string) error {
	ctx, concelFunc := initContext(c.env)
	defer concelFunc()
	collec := c.getCollection()
	filter := bson.M{"email": email}
	r, err := collec.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	fmt.Printf("[DELETE]number of record was deleted: %v\n", r.DeletedCount)
	return nil
}

func (c user) Create(entity *UserEntity) error {
	ctx, cancelFunc := initContext(c.env)
	defer cancelFunc()

	collec := c.getCollection()
	r, err := collec.InsertOne(ctx, entity)
	if err != nil {
		return err
	}
	if r.InsertedID == nil {
		return errors.New("err while insert data")
	}
	entity.ID = r.InsertedID.(primitive.ObjectID) //set up id to entity
	fmt.Printf("[INSERT]insert successfully, id: %v\n", entity.ID)
	return nil
}

func (c user) UpdateById(entity UserEntity) error {
	ctx, cancelFunc := initContext(c.env)
	defer cancelFunc()
	collec := c.getCollection()
	set := bson.M{"$set": bson.M{
		"info.first_name": "Kritchat2",
	}}
	r, err := collec.UpdateByID(ctx, entity.ID, set)
	if err != nil {
		return err
	}
	fmt.Printf("[UPDATE_BY_ID]number of record was updated: %v\n", r.MatchedCount)
	return nil
}

func (c user) UpdateByCondition(entity UserEntity) error {
	ctx, cancelFunc := initContext(c.env)
	defer cancelFunc()
	collec := c.getCollection()
	on := bson.M{"email": entity.Email} //update for user who have email id match with condition
	set := bson.M{"$set": bson.M{
		"info.last_name": "Rojanaphruk2",
	}}
	r, err := collec.UpdateOne(ctx, on, set)
	if err != nil {
		return err
	}
	fmt.Printf("[UPDATE_BY_CONDITION]number of record was updated: %v\n", r.MatchedCount)
	return nil
}

func (c user) FindOneById(id primitive.ObjectID) (*UserEntity, error) {
	ctx, cancelFunc := initContext(c.env)
	defer cancelFunc()
	collec := c.getCollection()
	filter := bson.M{"_id": id}
	r := collec.FindOne(ctx, filter)
	if r.Err() != nil {
		return nil, r.Err()
	}
	var result UserEntity
	err := r.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func (c user) getCollection() *mongo.Collection {
	return c.client.Database(c.DBName).Collection(c.collectionName)
}
