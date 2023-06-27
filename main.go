package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CatFactWorker struct {
	client *mongo.Client
}

func NewCatFactWorker(c *mongo.Client) *CatFactWorker {
	return &CatFactWorker{
		client: c,
	}
}

func (cfw *CatFactWorker) start() error {
	collection := cfw.client.Database("catfact").Collection("facts")
	// ticker.C channel will be sent a value every 2 seconds
	ticker := time.NewTicker(2 * time.Second)

	for {
		// http req to endpoint
		resp, err := http.Get("https://catfact.ninja/fact")
		if err != nil {
			return err
		}

		// declare new bson.M struct to store cat facts
		var catFact bson.M
		// decode json data in resp.body into catFact struct
		if err := json.NewDecoder(resp.Body).Decode(&catFact); err != nil {
			return err
		}

		_, err = collection.InsertOne(context.Background(), catFact)
		if err != nil {
			return err
		}
		<-ticker.C
	}
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	worker := NewCatFactWorker(client)
	worker.start()
}
