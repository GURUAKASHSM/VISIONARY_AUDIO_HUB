package models

// FeedBack Given
type Feedback struct {
	Token    string `json:"token" bson:"token"`
	Feedback string `json:"feedback" bson:"feedback"`
}

// FeedBack to Set In DB
type FeedbackDB struct {
    CustomerId         string `json:"customerid" bson:"customerid"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Feedback string `json:"feedback" bson:"feedback"`
}