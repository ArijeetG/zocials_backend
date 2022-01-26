package helpers

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// const uri = "mongodb+srv://users:Acaspera%40123@cluster0.0bnj6.mongodb.net/?retryWrites=true&w=majority"

func ConnectToDatabase() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://users:Acaspera%40123@cluster0.0bnj6.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	e := client.Connect(ctx)
	if e != nil {
		panic(e)
	}
	// defer client.Disconnect(ctx)
	return client
}
