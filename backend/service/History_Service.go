package service

import (
	"audiohub/config"
	"audiohub/constants"
	"audiohub/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func SaveToHistory(history models.History) (string, error) {
	var sitedata models.AdminPageData
	err := config.SiteDataCollection.FindOne(context.Background(), bson.M{}).Decode(&sitedata)
	if err != nil {
		log.Println(err)
		return "User Not Found", err
	}
	if history.Type == "PDF" {
		sitedata.PDFInputCount++
	} else if history.Type == "IMAGE" {
		sitedata.ImageInputCount++
	} else if history.Type == "TEXT" {
		sitedata.TextInputCount++
	}
	sitedata.HistoryCount++
	sitedata.AI_OutputCount++
	sitedata.OCR_Count++
	_, err = config.SiteDataCollection.ReplaceOne(context.Background(), bson.M{}, sitedata)
	if err != nil {
		log.Println(err)
		return "Error in Updateing Count", err
	}

	Id, err := ExtractCustomerID(history.CustomerId, constants.SecretKey)
	if err != nil {
		log.Println(err)
		return "Login Expired", err
	}

	var Customer models.Customer
	filter := bson.M{"customerid": Id}
	err = config.User_Collection.FindOne(context.Background(), filter).Decode(&Customer)
	if err != nil {
		log.Println(err)
		return "User Not Found", err
	}
	if history.Type == "PDF" {
		Customer.Pdfcount++
	} else if history.Type == "IMAGE" {
		Customer.Imagecount++
	} else if history.Type == "TEXT" {
		Customer.Textcount++
	}
	Customer.Totalcount++
	Customer.Ocrcount++
	Customer.AIcount++
	_, err = config.User_Collection.ReplaceOne(context.Background(), filter, Customer)
	if err != nil {
		log.Println(err)
		return "Error in Updateing Count", err
	}
	history.HistoryId = GenerateUniqueHistoryID()
	history.CustomerId = Id
	history.SavedTime = GetCurrentTimeAndDate()
	_, err = config.History_Collection.InsertOne(context.Background(), history)
	if err != nil {
		log.Println(err)
		return "Error in Inserting", err
	}
	return "Success", nil
}

func DisplayHistory(token models.Token) ([]models.History, string, error) {
	id, err := ExtractCustomerID(token.Token, constants.SecretKey)
	if err != nil {
		log.Print(err)
		return nil, "Login Expired", err
	}
	var History []models.History
	cursor, err := config.History_Collection.Find(context.Background(), bson.M{"customerid": id})
	if err != nil {
		log.Print(err)
		return nil, "History not Found", err
	}
	for cursor.Next(context.Background()) {
		var history models.History
		err := cursor.Decode(&history)
		if err != nil {
			return nil, "Error in Decode", err
		}
		History = append(History, history)
	}
	return History, "Success", nil
}

func DeteleHistory(details models.DeteleHistory) (string, error) {
	id, err := ExtractCustomerID(details.Token, constants.SecretKey)
	if err != nil {
		log.Println(err)
		return "Login Expired", err
	}
	filter1 := bson.M{"customerid": id}
	filter2 := bson.M{"historyid": details.Historyid}
	combinedfilter := bson.M{"$and": []bson.M{filter1, filter2}}
	_, err = config.History_Collection.DeleteOne(context.Background(), combinedfilter)
	if err != nil {
		log.Println(err)
		return "History Not found", err
	}
	filter := bson.M{}
	update := bson.M{
		"$inc": bson.M{"historycount": -1}, 
	}
	_, err = config.SiteDataCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
		return "Error in Updateing Site Data",err
	}
	return "Deleted Successfully",nil
}
func GetCurrentTimeAndDate() string {
	currentTime := time.Now()
	timeStr := currentTime.Format("3:04 PM")
	dateStr := currentTime.Format("2 Jan 2006")
	timeAndDateStr := fmt.Sprintf("%s, %s", timeStr, dateStr)
	return timeAndDateStr
}

func GenerateUniqueHistoryID() string {
	// Implement your logic to generate a unique customer ID (e.g., UUID, timestamp, etc.)
	// For example, you can use a combination of timestamp and random characters
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), GetRandomString(16))
}
