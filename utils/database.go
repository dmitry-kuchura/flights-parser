package utils

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection = mongo.Collection

func UpdateOne(collection *Collection, value string, data Flight) {
	filter := bson.D{{"number", value}}

	update := bson.D{
		{"$set", bson.D{
			{"boardStatus", data.BoardStatus},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func InsertMany(collection *Collection, list Flight) {
	flights := []interface{}{list}

	insertManyResult, err := collection.InsertMany(context.TODO(), flights)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}

func InsertOne(collection *Collection, flight Flight) {
	insertResult, err := collection.InsertOne(context.TODO(), flight)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func FindOne(collection *Collection, value string) (flight Flight, err error) {
	filter := bson.D{{"number", value}}

	err = collection.FindOne(context.TODO(), filter).Decode(&flight)

	return flight, err
}
