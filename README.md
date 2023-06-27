# Cat Facts API

- First simple project for practicing Go
- Sends Request to catfact.ninja/facts endpoint to retrieve random facts every 2 seconds, stores the fact in mongoDB if it doesn't exist yet.

## Concepts practiced

- pointers
- Goroutines and channels
- mongoDB with Go (bson, cursor etc.)
- lean/basic err handling
- http package
- time package (ticker)

# Set up

## Installing mongodb

### Installing mongodb with Docker

```
docker run --name some-mongo -p 27017:27017 -d mongo
```

### Go dependencies

```
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/bson
```

### Mongo Golang quickstart

```
client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
```
