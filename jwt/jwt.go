package returnjwt

import (
	"encoding/json"
	"fmt"
	"nank/app/structs"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var key = []byte("myseckey")

func GenerateToken(phone string, permision string, id string) string {
	expiration := time.Now().Add(5 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, structs.Tokenclaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
			Subject:   "hello",
		},
		UserId:     id,
		Phone:      phone,
		Permission: permision,
	})
	fmt.Printf("token: %v\n", token)
	// standartclaim joyte kenom xo claim ar ve darun json ved name surname done

	fmt.Printf("token: %v\n", token)
	stringtoken, err := token.SignedString(key)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	return stringtoken
}
func Decode(token *jwt.Token) (interface{}, error) {
	return key, nil
}

func Validate(token string) (structs.Tokenclaim, bool) {
	tokenress, err := jwt.ParseWithClaims(token, &structs.Tokenclaim{}, Decode)
	if err == nil {
		fmt.Printf("err: %v\n", err)
	}
	ValidKey := tokenress.Claims

	var shablon structs.Tokenclaim

	databyte, _ := json.Marshal(ValidKey)
	json.Unmarshal(databyte, &shablon)
	fmt.Printf("shablon: %v\n", shablon.Permission)

	return shablon, tokenress.Valid
}
