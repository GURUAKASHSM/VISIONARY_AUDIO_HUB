package models

// Store History
type History struct {
	HistoryId  string `json:"historyid" bson:"historyid"`
	SavedTime  string `json:"savedtime" bson:"savedtime"`
	Type       string `json:"type" bson:"type"`
	File       string `json:"file" bson:"file"`
	CustomerId string `json:"customerid" bson:"customerid"`
	OCR_Text   string `json:"ocrtext" bson:"ocrtext"`
	AI_Text    string `json:"aitext" bson:"aitext"`
	OCR_Audio  string `json:"ocraudio" bson:"ocraudio"`
	AI_Audio   string `json:"aiaudio" bson:"aiaudio"`
}

// Delete History
type DeteleHistory struct {
	Historyid string `json:"historyid" bson:"historyid"`
	Token     string `json:"token" bson:"token"`
}
