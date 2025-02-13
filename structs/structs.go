package structs

import "github.com/dgrijalva/jwt-go"

type Registration struct {
	Id          string `bson:"_id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Verify bool `json:"verify"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	SecrateCode        int    `json:"secratecode"` // Добавлено поле Code
	Cash int  `json:"cash"`
	Permission string `json:"permission"`
}
type UpdatePass struct {
	Password string `json:"password"`
	Newpassword string `json:"Newpassword"`
	Email string `json:"email"`
}

type Tokenclaim struct {
	jwt.StandardClaims
	UserId     string `bson:"_id"`
	Phone      string `json:"phone"`
	Permission string `json:"permission"`
}



type AddMathod struct{
	Id string  `bson:"_id"`
	Name string
	Logo string
	Description string
}



type Delete_Mathod struct{
	Id string `bson:"_id"`
	LogoId string
}
type Verify struct {
	Id   string `bson:"_id"`
	Code  int `json:"code"`
	Email string `json:"email"`
	User_Id string `json:"user_id"`
}
type UpdatePassword struct{
	Email string `json:"email"`
	Code int `json:"code"`
}
type CheckPassword struct {
	Id   string `bson:"_id"`
	Email string `json:"email"`
	Password string `json:"password"`
}


type CashStruct struct{
	Id   string `bson:"_id"`
	Cash int `json:"cash"`
	Phone string `json:"phone"`
	SenderNumber string `json:"senderNumber"`
	ReciverNumber string `json:"reciverNumber"`
	Sender_id string `json:"sender_id"`
	Reciver_id string `json:"reciver_id"`
}

type AddCash struct{
	Id string `bson:"_id"`
	Cash int `json:"cash"`
	Phone string `json:"phone"`
}
