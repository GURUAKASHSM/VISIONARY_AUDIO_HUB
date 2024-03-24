package models

// Sign Up Customer
type Customer struct {
	CustomerId         string `json:"customerid" bson:"customerid"`
	Email              string `json:"email" bson:"email"`
	Name               string `json:"name" bson:"name"`
	Phone_No           int    `json:"phonenumber" bson:"phonenumber"`
	Password           string `json:"password" bson:"password"`
	ConfirmPassword    string `json:"confirmpassword" bson:"confirmpassword"`
	IsEmailVerified    bool   `json:"isemailverified" bson:"isemailverified"`
	WrongInput         int    `json:"wronginput" bson:"wronginput"`
	VerificationString string `json:"verification" bson:"verification"`
	BlockedUser        bool   `json:"blockeduser" bson:"blockeduser"`
	Image              string `json:"image"bson:"image"`
	Date               string `json : "createddate" bson"createddate"`
	Pdfcount           int64  `json : "pdfcount" bson"pdfcount"`
	Imagecount         int64  `json:"imagecount" bson"imagecount"`
	Textcount          int64  `json:"textcount" bson"textcount"`
	Totalcount         int64  `json:"totalcount" bson"totalcount"`
	Ocrcount           int64  `json:"ocrcount" bson"ocrcount"`
	AIcount            int64  `json:"aicount" bson"aicount"`
}

// Customer Login
type Login struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type VerifyEmail struct {
	Email              string `json:"email" bson:"email"`
	VerificationString string `json:"verification" bson:"verification"`
}
type ForgetPassword struct {
	Email string `json:"email" bson:"email"`
}
type PasswordChange struct {
	Email              string `json:"email" bson:"email"`
	VerificationString string `json:"verification" bson:"verification"`
	Password           string `json:"password" bson:"password"`
	ConfirmPassword    string `json:"confirmpassword" bson:"confirmpassword"`
}
