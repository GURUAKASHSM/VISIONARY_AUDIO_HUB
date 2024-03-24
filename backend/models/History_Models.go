package models

//Store History
type History struct{
	CustomerId         string `json:"customerid" bson:"customerid"`
	OCR_Text string `json:"ocrtext" bson:"ocrtext"`
	AI_Text string `json:"aitext" bson:"aitext"`
	OCR_Audio string `json"ocraudio" bson:"ocraudio"`
	AI_Audio string `json:"aiaudio" bson:"aiaudio"`
}