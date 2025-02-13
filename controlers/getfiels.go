package controlers

import (
	"fmt"
	env "nank/app/Env"
	mongoconnect "nank/app/MongoConnect"
	"nank/app/structs"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)
func Verify(c *gin.Context) {
	email := c.Query("Email")
	code := c.Query("SecrateCode")

	if email == "" || code == "" {
		c.JSON(400, gin.H{"error": "Email and SecrateCode are required"})
		return
	}
	client, ctx := mongoconnect.DBConnection()
	defer client.Disconnect(ctx)
	str :=code
	num, error := strconv.ParseFloat(str, 64) 
	if error != nil {
		fmt.Println("Ошибка:", error)
	}
	collection := client.Database(env.Data_Name).Collection("code")
	var codeDoc structs.Registration
	err := collection.FindOne(ctx, bson.M{
		"email":       email,
		"secrateCode": num,
	}).Decode(&codeDoc)

	if err != nil {
		c.JSON(404, gin.H{"error": "Invalid code or email"})
		return
	}
	collection2 := client.Database(env.Data_Name).Collection("Users")
	updateResult, err := collection2.UpdateOne(ctx,
		bson.M{"email": email},
		bson.M{"$set": bson.M{"verify": true}},
	)

	if err != nil || updateResult.ModifiedCount == 0 {
		c.JSON(500, gin.H{"error": "Failed to verify user"})
		return
	}
	c.JSON(200, gin.H{"message": "Verification successful"})
}