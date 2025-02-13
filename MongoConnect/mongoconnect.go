package mongoconnect

import (
	"context"
	"fmt"
	"nank/app/Env"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBConnection() (*mongo.Client, context.Context) {
	url := options.Client().ApplyURI(env.Data_Ip)
	NewCtx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	Client, err := mongo.Connect(NewCtx, url)
	if err != nil {
		fmt.Printf("errors: %v\n", err)
	}
	return Client, NewCtx
}
