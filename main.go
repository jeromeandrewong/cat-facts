package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	client *mongo.Client
}

func NewServer(c *mongo.Client) *Server {
	return &Server{
		client: c,
	}
}

func (s *Server) handleGetAllFacts(w http.ResponseWriter, r *http.Request) {
	// get reference to "facts" collection in "catfact" database
	collection := s.client.Database("catfact").Collection("facts")

	// query that matches all docs in collection
	query := bson.M{}
	cursor, err := collection.Find(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	results := []bson.M{} // slice of bson.M structs to store results
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)                       // set status code to 200
	w.Header().Add("Content-Type", "application/json") // set content type to json
	json.NewEncoder(w).Encode(results)                 // encode results into json and write to response writer

}

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
	go worker.start()

	server := NewServer(client)
	http.HandleFunc("/facts", server.handleGetAllFacts)

	http.ListenAndServe(":3000", nil)
}
