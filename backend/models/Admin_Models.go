package models

// To Delete Data
type Delete struct {
	Token string `json:"token" bson:"token"`
	
	Id    string `json:"id" bson:"id"`
}

// To Upadte Feild
type Update struct {
	Token  string `json:"token" bson:"token"`
	Id     string `json:"id" bson:"id"`
	Field      string `json:"field" bson:"field"`
	New_Value  string `json:"newvalue" bson:"newvalue"`
}

// Admin Signup Data
type AdminData struct {
	AdminID    string `json:"adminid" bson:"adminid"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	TOTP       string `json:"totp" bson:"totp"`
	IP_Address string `json:"ip" bson:"ip"`
	SecretKey  string `json:"secretkey" bson:"secretkey"`
	WrongInput int    `json:"wronginput" bson:"wronginput"`
	Token      string `json:"token" bson:"token"`
}

// Admin Sign in data
type AdminSignup struct {
	AdminName       string `json:"name" bson:"name"`
	Email           string `bson:"email" json:"email"`
	Password        string `bson:"password" json:"password"`
	ConfirmPassword string `bson:"confirmpassword" json:"confirmpassword"`
	IP              string `bson:"ip" json:"ip"`
}

// Data Needed for Admin Page
type AdminPageData struct {
	UserCount        int64 `json:"usercount" bson:"usercount"`
	HistoryCount      int64 `json:"historycount" bson:"historycount"`
	AI_OutputCount     int64 `json:"aicount" bson:"aicount"`
	OCR_Count       int64 `json:"ocrcount" bson:"ocrcount"`
	ImageInputCount int64 `json:"imageinputcount" bson:"imageinputcount"`
	TextInputCount int64 `json:"textinputcount" bson:"textinputcount"`
	PDFInputCount int64 `json:"pdfinputcount" bson:"pdfinputcount"`
}

// Get Every Single Data
type Getdata struct {
	Id         string `json:"id" bson:"id"`
	Collection string `json:"collection" bson:"collection"`
}

// Single Data Returing Structure
type ReturnData struct {
	CustomerId         string `json:"customerid" bson:"customerid"`
	Name               string `json:"name" bson:"name"`
	IsEmailVerified    bool   `json:"isemailverified" bson:"isemailverified"`
	WrongInput         int    `json:"wronginput" bson:"wronginput"`
	VerificationString string `json:"verification" bson:"verification"`
	BlockedUser        bool   `json:"blockeduser" bson:"blockeduser"`	
	Image           string `json:"image" bson:"image"`
	Email           string `json:"email" bson:"email"`	
	Phone_No        int    `json:"phonenumber" bson:"phonenumber"`
	Password        string `json:"password" bson:"password"`
	ConfirmPassword string `json:"confirmpassword" bson:"confirmpassword"`
}

// Block User
type Block struct {
	Email      string `json:"email" bson:"email"`
	Collection string `json:"collection" bson:"collection"`
}

// ShutDown
type ShutDown struct {
	Token    string `json:"token" bson:"token"`
	Password string `json:"password" bson:"password"`
}