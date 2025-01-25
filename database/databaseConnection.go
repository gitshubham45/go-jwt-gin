package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatal("Error loading .env file")
	}

	MongoDbUri := os.Getenv("MONGODB_URL")

	ctx , cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() 

	client , err := mongo.Connect(ctx,options.Client().ApplyURI(MongoDbUri))
	if err != nil{
		log.Fatal("Error connecting to mongoDB ",err)
	}

	fmt.Println("Mongo DB connected")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client , collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}