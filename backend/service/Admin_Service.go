package service

import (
	"audiohub/config"
	"audiohub/constants"
	"audiohub/models"
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Admin Login
func AdminLoginCheck(login *models.AdminData) (string, int) {

	var correctdata models.AdminData
	filter := bson.M{"email": login.Email}
	err := config.Admin_Collection.FindOne(context.Background(), filter).Decode(&correctdata)
	if err != nil {
		return "", 0
	}
	if correctdata.WrongInput == 4 {
		return "", 1
	}
	if correctdata.IP_Address == login.IP_Address {
		return "", 3
	}
	if correctdata.Password != login.Password {
		correctdata.WrongInput++
		update := bson.M{"$set": bson.M{"wronginput": correctdata.WrongInput}}
		config.Admin_Collection.UpdateOne(context.Background(), filter, update)
		return "", 2
	}

	if !ValidateOTP(login.TOTP, correctdata.SecretKey) {
		correctdata.WrongInput++
		update := bson.M{"$set": bson.M{"wronginput": correctdata.WrongInput}}
		config.Admin_Collection.UpdateOne(context.Background(), filter, update)
		return "", 4
	}

	token, err := CreateToken(correctdata.Email, correctdata.AdminID)
	if err != nil {
		return "", 5
	}

	log.Println(token)
	update := bson.M{"$set": bson.M{"wronginput": 0}}
	config.Admin_Collection.UpdateOne(context.Background(), filter, update)
	return token, 5

}

// TO get all Customer
func GetallCustomerdata() ([]models.Customer, string, error) {
	filter := bson.D{}
	cursor, err := config.User_Collection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(context.Background())
	var Profiles []models.Customer
	for cursor.Next(context.Background()) {
		var profile models.Customer
		err := cursor.Decode(&profile)
		if err != nil {
			return nil, "Error in Decode", err
		}
		Profiles = append(Profiles, profile)
	}
	return Profiles, "Success", nil
}

// Update Any Data
func Update(update models.Update) bool {

	filter := bson.M{"selleremail": update.Id}
	update1 := bson.M{"$set": bson.M{update.Field: update.New_Value}}
	options := options.Update()
	_, err := config.User_Collection.UpdateOne(context.TODO(), filter, update1, options)
	if err != nil {
		return false
	}
	go SendEditDataNotification(update.Id, update.Field, update.New_Value)
	return true

}

// Delete Any data
func Delete(delete models.Delete) bool {

	filter := bson.M{"email": delete.Id}
	_, err := config.User_Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

// Get Dataneed for Admin
func AdminNeededData() (*models.AdminPageData, string, error) {
	var adminpagedata models.AdminPageData
	err := config.SiteDataCollection.FindOne(context.Background(), bson.M{}).Decode(&adminpagedata)
	if err != nil {
		log.Println("Error  inFinding")
		return nil, "No DataFound", err
	}
	return &adminpagedata, "Success", nil
}

// Create Admin
func CreateAdmin(admin models.AdminSignup) (string, string) {
	filter := bson.M{"email": admin.Email}

	result := config.Admin_Collection.FindOne(context.Background(), filter)
	if result.Err() == nil {
		return "User Already Exists", ""
	}
	if result.Err() != nil && result.Err() != mongo.ErrNoDocuments {
		return "Error in Query: " + result.Err().Error(), ""
	}
	key, err := GenerateSecret()
	if err != nil {
		return "Error In Generating TOTP", ""
	}
	var AdminData models.AdminData
	AdminData.Email = admin.Email
	AdminData.Password = admin.Password
	AdminData.IP_Address = admin.IP
	AdminData.SecretKey = key
	AdminData.Token = ""
	AdminData.WrongInput = 0
	AdminData.AdminID = GenerateUniqueAdminID()
	_, err = config.Admin_Collection.InsertOne(context.Background(), AdminData)
	if err != nil {
		return "Error in Creating: " + err.Error(), ""
	}
	go SendAdminInvitation(admin.Email, admin.AdminName, admin.Password, "https://anon.up.railway.app/admin/", admin.IP, key)
	return "Created Successfully", key
}

// Get Single Data
func GetData(data models.Getdata) (*models.ReturnData, error) {
	var returndata models.ReturnData

	if data.Collection == "customer" {
		log.Println("In customer")
		var profile models.Customer
		filter := bson.M{"email": data.Id}
		err := config.User_Collection.FindOne(context.Background(), filter).Decode(&profile)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		returndata.Name = profile.Name
		returndata.CustomerId = profile.CustomerId

		returndata.Email = profile.Email
		returndata.Phone_No = profile.Phone_No
		returndata.Password = profile.Password
		returndata.IsEmailVerified = profile.IsEmailVerified
		returndata.BlockedUser = profile.BlockedUser
		returndata.WrongInput = profile.WrongInput
		return &returndata, nil

	}
	return nil, nil

}

// Block User & Admin
func Block(data models.Block) (string, error) {
	if data.Collection == "customer" {
		var customer models.Customer
		filter := bson.M{"email": data.Email}
		err := config.User_Collection.FindOne(context.Background(), filter).Decode(&customer)
		if err != nil {
			log.Println(err)
			return "No result Found", err
		}
		message := ""
		if customer.BlockedUser {
			customer.BlockedUser = false
			message = "Customer has been Unblocked"
		} else {
			customer.BlockedUser = true
			message = "Customer has been Blocked"
		}
		update := bson.M{"$set": bson.M{"blockeduser": customer.BlockedUser}}
		_, err = config.User_Collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			log.Println(err)
			return "Can't Update Data", err
		}
		go SendBlockingNotification(customer.Email, customer.Name, "Due to improper behaviour")
		return message, nil
	}
	return "Invalid Collection", nil
}

// ShutDown
func ShutDown(token models.ShutDown) (string, error) {
	if token.Password != constants.ShutDownKey {
		return "Key Mismatch", nil
	}
	id, err := ExtractCustomerID(token.Token, constants.SecretKey)
	if err != nil {
		log.Println("Login Exp")
		return "Login Expired", err
	}
	var admin models.AdminData
	filter := bson.M{"adminid": id}
	err = config.Admin_Collection.FindOne(context.Background(), filter).Decode(&admin)
	if err != nil {
		return "Login as Admin", err
	}

	if id != admin.AdminID {
		return "Login as Admin", err
	}

	shutdownComplete := make(chan bool)

	go func() {
		ShutDownExe()
		shutdownComplete <- true
	}()

	return "Shutdown initiated", nil
}

func ShutDownExe() {
	time.Sleep(3 * time.Second)
	os.Exit(0)
}

// Clear DataBase
func ClearDB(details models.Getdata) (string, error) {
	id, err := ExtractCustomerID(details.Id, constants.SecretKey)
	if err != nil {
		return "Login Expired", err
	}
	var admin models.AdminData
	filter := bson.M{"adminid": id}
	err = config.Admin_Collection.FindOne(context.Background(), filter).Decode(&admin)
	if err != nil {
		return "Data not Found", err
	}
	if admin.Email == "" {
		return "Data not Found", nil
	}
	result, err := DeleteDBCollection(details.Collection)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Delete colletion
func DeleteDBCollection(collection string) (string, error) {
	if collection == "all" {
		err := config.Admin_Collection.Drop(context.Background())
		if err != nil {
			return "Error in Delting Admin Collection", err
		}

		err = config.User_Collection.Drop(context.Background())
		if err != nil {
			return "Error in Delting Customers Collection", err
		}

		err = config.History_Collection.Drop(context.Background())
		if err != nil {
			return "Error in Delting Inventory Collection", err
		}

		return "All Database Deleted Successfully", nil

	} else if collection == "userall" {
		err := config.User_Collection.Drop(context.Background())
		if err != nil {
			return "Error in Delting Orders Collection", err
		}
		err = config.History_Collection.Drop(context.Background())
		if err != nil {
			return "Error in Delting History Collection", err
		}

		return "User Related Database Deleted Successfully", nil
	} else if collection == "user" {
		err := config.User_Collection.Drop(context.Background())
		if err != nil {
			return "Error in Delting Admin Collection", err
		}

		return "User  Database Deleted Successfully", nil
	} else if collection == "admin" {
		err := config.Admin_Collection.Drop(context.Background())
		if err != nil {
			return "Error in Delting Sellers Collection", err
		}
		return "Admin Database Deleted Successfully", nil
	} else if collection == "history" {
		err := config.History_Collection.Drop(context.Background())
		if err != nil {
			return "Error in Delting Inventory Collection", err
		}
		return "History Database Deleted Successfully", nil
	}
	return "Collection Not Found", nil

}
