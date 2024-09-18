package initializers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SHRYNSH-NETAM/Sudoku_Backend/models"
	"github.com/redis/go-redis/v9"
)

var (
	ctx    = context.Background()
	client *redis.Client // Declare client at package level
)

func Connect2Redis(){
	uri := os.Getenv("REDIS_URI")
	if uri==""{
		log.Fatal("Set your 'REDIS_URI' environment variable in .env file")
	}

	fmt.Println("Connecting to Redis...")

	opt, err := redis.ParseURL(uri)
	if err != nil {
		panic(err)
	}

	client = redis.NewClient(opt) //Initializing the client variable

	_, err = client.Ping(context.Background()).Result()
	if err!=nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Successfully connected and pinged to Redis")
}

func Add2Redis(key string, result models.Sudokugrid) error{
	if client == nil {
		fmt.Println("Redis client is not initialized")
        return nil
    }
	
	jsonString, err := json.Marshal(models.Sudokugrid{
		SolvedGrid: result.SolvedGrid,
		UnSolvedGrid: result.UnSolvedGrid,
		Time: result.Time,
	})
	if err!=nil {
		return err
	}

	// key := fmt.Sprintf("Result:%s",result.ID)
	err = client.Set(ctx, key, jsonString, 1*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func Get2Redis(key string) (models.Sudokugrid, error){
	var Cacheresult models.Sudokugrid
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return Cacheresult,err
	}
	err = json.Unmarshal([]byte(val), &Cacheresult)
    if err != nil {
        return Cacheresult, err
    }
	return Cacheresult,nil
}