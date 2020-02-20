package utils

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection = mongo.Collection

func UpdateOne(collection *Collection, value *string, data *Flight) {
	filter := bson.D{{"Number", value}}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func InsertMany(collection *Collection, list *Flight) {
	flights := []interface{}{list}

	insertManyResult, err := collection.InsertMany(context.TODO(), flights)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}

func InsertOne(collection *Collection, flight *Flight) {
	insertResult, err := collection.InsertOne(context.TODO(), flight)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func FindOne(collection *Collection, value *string) {
	var result Flight

	filter := bson.D{{"Number", value}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)
}
