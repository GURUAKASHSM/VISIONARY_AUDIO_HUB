package service

import (
	"audiohub/config"
	"audiohub/constants"
	"audiohub/models"
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create Customer
func CreateCustomer(profile models.Customer) int {
	name := profile.Name
	regexPattern := "^[a-zA-Z ]+$"
	regexpObject := regexp.MustCompile(regexPattern)
	match := regexpObject.FindString(name)

	if match == "" {
		return 4
	}

	phonecount := countdigits(profile.Phone_No)

	if phonecount != 10 {
		return 5
	}

	if profile.Password != profile.ConfirmPassword {
		return 3
	}

	filter := bson.M{"email": profile.Email}
	existingCustomer := config.User_Collection.FindOne(context.Background(), filter)

	// Check if the document already exists
	if existingCustomer.Err() == nil {
		var existingProfile models.Customer // Replace with the actual type of your document
		if err := existingCustomer.Decode(&existingProfile); err != nil {
			log.Println(err)
			return 0
		}

		// Check if the email is verified
		if existingProfile.IsEmailVerified {
			return 2
		} else {
			profile.CustomerId = GenerateUniqueUserID()
			profile.IsEmailVerified = false
			profile.WrongInput = 0
			profile.VerificationString = GenerateOTP(6)
			if profile.Image == "" {
				profile.Image = "No Image"
			}
			profile.AIcount = 0
			profile.Date = time.Now().Format("02 Jan 2006")
			profile.Imagecount = 0
			profile.Ocrcount = 0
			profile.Pdfcount = 0
			profile.Textcount = 0
			profile.Totalcount = 0
			_, updateErr := config.User_Collection.ReplaceOne(context.Background(), filter, profile)
			if updateErr != nil {
				return 0
			}
			go SendEmailforVerification(profile.Email, profile.VerificationString, profile.Name)

			return 1
		}
	} else if existingCustomer.Err() == mongo.ErrNoDocuments {

		profile.CustomerId = GenerateUniqueUserID()

		profile.IsEmailVerified = false
		profile.BlockedUser = false
		profile.WrongInput = 0
		profile.VerificationString = GenerateOTP(6)
		if profile.Image == "" {
			profile.Image = "No Image"
		}
		profile.AIcount = 0
		profile.Date = time.Now().Format("02 Jan 2006")
		profile.Imagecount = 0
		profile.Ocrcount = 0
		profile.Pdfcount = 0
		profile.Textcount = 0
		profile.Totalcount = 0

		inserted, insertErr := config.User_Collection.InsertOne(context.Background(), profile)
		if insertErr != nil {
			log.Println(insertErr)
			return 0
		}
		go SendEmailforVerification(profile.Email, profile.VerificationString, profile.Name)
		fmt.Println("Inserted", inserted.InsertedID)
		return 1
	} else {
		log.Println(existingCustomer.Err())
		return 0
	}

}

// Email Verification for Customer
func EmailVerification(data models.VerifyEmail) (string, error) {
	filter := bson.M{"email": data.Email}
	var customer models.Customer
	err := config.User_Collection.FindOne(context.Background(), filter).Decode(&customer)
	if err != nil {
		return "", err
	}
	if customer.IsEmailVerified {
		return "Email already verified", nil
	}

	if customer.VerificationString == data.VerificationString {
		update := bson.M{"$set": bson.M{"isemailverified": true}}
		filter := bson.M{"email": data.Email}

		_, err := config.User_Collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return "", err
		}
		go SendThankYouEmail(customer.Email, customer.Name)
		go IncrementUserCount()
		return "Signup Successful", nil
	} else {
		return "Wrong OTP", nil
	}

}

// Incrementuser Count In SiteData DB
func IncrementUserCount() {
	update := bson.M{
		"$inc": bson.M{"usercount": 1},
	}
	_, err := config.SiteDataCollection.UpdateOne(context.Background(), bson.M{}, update)
	if err != nil {
		log.Println(err)
	}
}

// Login Customer
func Login(details models.Login) (string, int, error) {
	var customer models.Customer

	filter := bson.M{"email": details.Email}
	err := config.User_Collection.FindOne(context.Background(), filter).Decode(&customer)
	if err != nil {
		return "User not found", 0, err
	}
	if customer.WrongInput == 10 {
		return "Too many no of try", 0, nil
	}
	if !customer.IsEmailVerified {
		return "Please verify your email", 0, nil
	}
	if customer.BlockedUser {
		return "Your ID has been Blocked", 0, nil
	}
	if customer.Password != details.Password {
		customer.WrongInput++
		update := bson.M{"$set": bson.M{"wronginput": customer.WrongInput}}
		config.User_Collection.UpdateOne(context.Background(), filter, update)
		return "Wrong Password", 0, nil
	}

	token, err := CreateToken(customer.Email, customer.CustomerId)
	if err != nil {
		return "Internal server error", 0, nil

	}

	return token, 1, nil
}
func GenerateUniqueUserID() string {
	// Implement your logic to generate a unique customer ID (e.g., UUID, timestamp, etc.)
	// For example, you can use a combination of timestamp and random characters
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), GetRandomString(6))
}

func SendEmailForForgotPassword(email models.ForgetPassword) (string, error) {

	filter := bson.M{"email": email.Email}
	var customer models.Customer
	err := config.User_Collection.FindOne(context.Background(), filter).Decode(&customer)
	if err != nil {
		return "Email not Found", err
	}
	if !customer.IsEmailVerified {
		return "Email not verified", nil
	}
	customer.WrongInput = 0
	customer.VerificationString = GenerateOTP(6)
	_, updateErr := config.User_Collection.ReplaceOne(context.Background(), filter, customer)
	if updateErr != nil {
		return "Invalid data", nil
	}
	go SendEmailforVerification(customer.Email, customer.VerificationString, customer.Name)
	return "Email sent successfully", nil

}
func ChangePassword(data models.PasswordChange) (string, error) {
	filter := bson.M{"email": data.Email}
	var customer models.Customer
	err := config.User_Collection.FindOne(context.Background(), filter).Decode(&customer)
	if err != nil {
		return "Email Not Found", err
	}
	if data.Password != data.ConfirmPassword {
		return "Password Mismatch", nil
	}
	if customer.VerificationString == data.VerificationString {
		update := bson.M{
			"$set": bson.M{
				"password":        data.Password,
				"confirmpassword": data.ConfirmPassword,
				"wronginput":      0,
			},
		}

		filter := bson.M{"email": data.Email}

		_, err := config.User_Collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return "", err
		}
	}
	return "Password changed successfully", nil
}
func UserDetails(token models.Token) (*models.Customer, string, error) {
	id, err := ExtractCustomerID(token.Token, constants.SecretKey)
	if err != nil {
		log.Print(err)
		return nil, "Login Expired", err
	}
	filter := bson.M{"customerid": id}
	var customer models.Customer
	err = config.User_Collection.FindOne(context.Background(), filter).Decode(&customer)
	if err != nil {
		log.Print(err)
		return nil, "no user", err
	}

	return &customer, "success", nil
}
