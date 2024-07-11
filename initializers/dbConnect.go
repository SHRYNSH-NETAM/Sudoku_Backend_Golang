package initializers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

func Connect2DB(){
	fmt.Println("Connecting to DB....")

	uri := os.Getenv("MONGO_URI")
	if uri=="" {
		log.Fatal("Set your 'MONGODB_URI' environment variable in .env file")
	}

	client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(uri))
	if err!=nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err!=nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}

	fmt.Println("Successfully connected and pinged to MONGODB.")
}

func FindData(loginData models.Fuser) *models.User{

	uri := os.Getenv("MONGO_URI")
	if uri=="" {
		log.Fatal("Set your 'MONGODB_URI' environment variable in .env file")
	}

	client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(uri))
	if err!=nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	coll := client.Database("test").Collection("users")
	var result models.User

	if err := coll.FindOne(context.TODO(), loginData).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.Panic(err)
	}

	return &result
}

func AddData(loginData models.User) bool {

	uri := os.Getenv("MONGO_URI")
	if uri=="" {
		log.Fatal("Set your 'MONGODB_URI' environment variable in .env file")
	}
	client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(uri))
	if err!=nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	coll := client.Database("test").Collection("users")
	_, err = coll.InsertOne(context.TODO(), loginData)
	return err == nil
}

func DeleteData(loginData models.Fuser) bool {

	uri := os.Getenv("MONGO_URI")
	if uri=="" {
		log.Fatal("Set your 'MONGODB_URI' environment variable in .env file")
	}
	client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(uri))
	if err!=nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	coll := client.Database("test").Collection("users")
	_, error := coll.DeleteOne(context.TODO(),loginData)
	return error==nil
}

func UpdateData(loginData models.Fuser, updatedData models.User) bool {

	uri := os.Getenv("MONGO_URI")
	if uri=="" {
		log.Fatal("Set your 'MONGODB_URI' environment variable in .env file")
	}
	
	client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(uri))
	if err!=nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	coll := client.Database("test").Collection("users")

	update := bson.M{
		"$set": updatedData,
	}

	_, error := coll.UpdateOne(context.TODO(), loginData, update)
	if error != nil {
		log.Println("Error updating document:", err)
		return false
	}

	return true
}