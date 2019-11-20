/**
	The responsibility of this service are as below

	1. Consume message from topic 'search-result' and persist in a Mongo database
	2. Notify availablity of complete search result in 'result-available' topic after sucessfully persisting
	3. Provide a rest API to GET search results for a set of uuid

	Note: Point 3 aboive can be implemented as a separate independently scalable microservice.package manager
	Currently implementing it as part of result manager service for sake of simplicity. Design wise it should
	be a separate service
**/
package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gopkg.in/mgo.v2/bson"
)

func main() {

	// Prepare Kafka Consumer
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	db, err := connectDB()
	c := make(chan interface{})

	go storeData(db.Collection("search-result"), c)

	// Subscribe to 'searchengine-result' Kafka topic
	kafkaConsumer.SubscribeTopics([]string{"searchengine-result"}, nil)

	// Run an infinite loop and listen for incoming messages
	// When new message comes process the message by caching in redis
	// until all three search engine results are recieved for a serach id
	for {
		msg, err := kafkaConsumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s. \n", msg.TopicPartition)
			bsonMsg, err := processMessage(msg.Value)
			if err != nil {
				fmt.Println("Error while unbarshalling 'search-result'")
			} else {
				c <- bsonMsg
			}
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func processMessage(message []byte) (interface{}, error) {
	var bdoc interface{}
	err := bson.UnmarshalJSON(message, &bdoc)
	if err != nil {
		fmt.Println("Error unmarshallng search result from topic search-result: ", err)
		return nil, err
	}
	return bdoc, nil
}

// storeData stores a search result in to the mongo persistance
func storeData(collection *mongo.Collection, c chan interface{}) {
	data := <-c
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	insertResult, err := collection.InsertOne(ctx, data)
	defer cancel()

	if err != nil {
		fmt.Println("Error while inserting data: ", err)
	} else {
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	}
}

// connectToDB starts a new database connection and returns a reference to it
func connectDB() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	options := options.Client().ApplyURI("mongodb://mongo-server:27017")
	options.SetMaxPoolSize(3)
	client, err := mongo.Connect(ctx, options)
	if err != nil {
		fmt.Println("Error while connecting to mongo server.")
		return nil, err
	}
	return client.Database("search-db"), nil
}
