package config

import (
	"audiohub/constants"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var User_Collection *mongo.Collection
var History_Collection *mongo.Collection
var Admin_Collection *mongo.Collection
var SiteDataCollection *mongo.Collection
var FeedbackCollection *mongo.Collection

func init() {
	clientoption := options.Client().ApplyURI(constants.Connectstring)
	client, err := mongo.Connect(context.TODO(), clientoption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDb sucessfully connected")
	User_Collection = client.Database(constants.DB_Name).Collection(constants.User_Collection)
	History_Collection = client.Database(constants.DB_Name).Collection(constants.History_Collection)
	Admin_Collection = client.Database(constants.DB_Name).Collection(constants.Admin_Collection)
	SiteDataCollection = client.Database(constants.DB_Name).Collection(constants.SiteData_Collection)
	FeedbackCollection = client.Database(constants.DB_Name).Collection(constants.Feedback_Collection)

	fmt.Println("All Collection Connected")
}
