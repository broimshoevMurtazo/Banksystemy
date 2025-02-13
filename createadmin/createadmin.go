package createadmin

import (
	"context"
	"fmt"
	env "nank/app/Env"
	hashedpasswod "nank/app/HashedPassword"
	mongoconnect "nank/app/MongoConnect"
	"nank/app/halpers"
	"nank/app/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Createadmin() {
	client, ctx := mongoconnect.DBConnection()
	Connections := client.Database(env.Data_Name).Collection("Users")
	// if admin exist not insert
	found := Connections.FindOne(ctx, bson.M{
		"permission": "Admin",
	})
	var shablon structs.Registration
	found.Decode(&shablon)

	if shablon.Permission == "Admin"  {
		fmt.Println("admin is exsist")
	} else {
		code := halpers.RandomSixDigit()
		ID := primitive.NewObjectID().Hex()
		Password :="Admin"
		Hashed,err := hashedpasswod.HashPassword(Password)
		if err!=nil{
			fmt.Printf("err: %v\n", err)
		}
		Connections.InsertOne(context.Background(), structs.Registration{
			Id: ID,
			Name: "Admin",
			Surname: "Admin",
			Phone:      "006116688",
			Password:  Hashed ,
			Permission: "Admin",
			Email:      "murtazobroimshoevm4333@gmail1.com",
			SecrateCode:code,
		})
	}
}



