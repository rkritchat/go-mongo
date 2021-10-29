package main

import (
	"fmt"
	"go-mongo/internal/collection"
	"go-mongo/internal/config"
)

func main() {
	cfg := config.InitConfig()
	defer cfg.Free()

	//init collection
	collection.CollectionMigration(cfg.Client, cfg.Env)

	userCollection := collection.NewUser(cfg.Client, cfg.Env)
	email := "rkritchat@gmail.com"
	//delete user if exist
	err := userCollection.DeleteByEmail(email)
	if err != nil {
		panic(err)
	}

	//create user
	userEntity := collection.UserEntity{
		Email: email,
		Credential: collection.UserCredential{
			Username:     "rkrichat",
			Password:     "1234",
			HintQuestion: "What is my cat name?",
			HintAwnser:   "Uefa",
		},
		Info: collection.UserInfo{
			Firstname: "Kritchat",
			Lastname:  "Rojanaphruk",
			Age:       28,
			Addreess:  "Thailand",
		},
	}
	err = userCollection.Create(&userEntity)
	if err != nil {
		panic(err)
	}

	//update first_name
	err = userCollection.UpdateById(userEntity)
	if err != nil {
		panic(err)
	}

	//update last_name
	err = userCollection.UpdateByCondition(userEntity)
	if err != nil {
		panic(err)
	}

	//find user by id
	fmt.Printf("start find user by id: %v\n", userEntity.ID)
	r, err := userCollection.FindOneById(userEntity.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\nuser: %#v\n", r)

	fmt.Println("connect db successfully")
}
