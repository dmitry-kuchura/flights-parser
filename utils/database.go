package utils

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateOne(collection *mongo.Collection, value string, data Flight) {
	filter := bson.D{{"number", value}}

	update := bson.D{
		{"$set", bson.D{
			{"boardstatus", data.BoardStatus},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func InsertOne(collection *mongo.Collection, flight Flight) {
	insertResult, err := collection.InsertOne(context.TODO(), flight)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func FindOne(collection *mongo.Collection, value string) (flight Flight, err error) {
	filter := bson.D{{"number", value}}

	err = collection.FindOne(context.TODO(), filter).Decode(&flight)

	return flight, err
}
