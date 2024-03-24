package service

import (
	"audiohub/config"
	"audiohub/constants"
	"audiohub/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// Instert FeedBack
func InstertFeedback(feedback models.Feedback) (string, error) {
	id, err := ExtractCustomerID(feedback.Token, constants.SecretKey)
	if err != nil {
		return "Login Expired", err
	}
	var user models.Customer
	filter := bson.M{"customerid": id}
	err = config.User_Collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return "Details not Found", nil
	}

	var FeedbackDB models.FeedbackDB
	FeedbackDB.Email = user.Email
	FeedbackDB.Name = user.Name
	FeedbackDB.Feedback = feedback.Feedback

	_, err = config.FeedbackCollection.InsertOne(context.Background(), FeedbackDB)
	if err != nil {
		log.Println(err)
		return "Error in Inserting", err
	}

	return "FeedBack Submited Successfully", nil
}

// Delete FeedBack
func Deletefeedback(delete models.FeedbackDB) int32 {
	filter1 := bson.M{"email": delete.Email}
	filter2 := bson.M{"feedback": delete.Feedback}
	combinedFilter := bson.M{
		"$and": []bson.M{filter1, filter2},
	}
	_, err := config.FeedbackCollection.DeleteMany(context.Background(), combinedFilter)
	if err != nil {
		return 0
	}
	return 1

}

// Get all Feedbacks
func GetFeedBacks() []models.Feedback {
	filter := bson.M{}
	cursor, err := config.FeedbackCollection.Find(context.Background(), filter)
	var Feedback []models.Feedback
	if err != nil {
		log.Println(err)
	}

	for cursor.Next(context.Background()) {
		var feedback models.Feedback
		err := cursor.Decode(&feedback)
		if err != nil {
			log.Println(err)
		}
		Feedback = append(Feedback, feedback)
	}
	return Feedback
}
