package controlers

import (
	"fmt"
	env "nank/app/Env"
	mongoconnect "nank/app/MongoConnect"
	returnjwt "nank/app/jwt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteMathod(c *gin.Context) {
	var cookidata, cookieerror = c.Request.Cookie(env.Data_Name)
	if cookieerror != nil {
		fmt.Printf("cookieerror: %v\n", cookieerror)
		c.JSON(404, "error Not Cookie found")
		fmt.Printf("cookidata: %v\n", cookidata)
	} else {
		SecretKeyData, isvalid := returnjwt.Validate(cookidata.Value)
		if SecretKeyData.Permission != "Admin" && isvalid {
			c.JSON(404, "error")
		} else {
			ids := c.Request.URL.Query().Get("id")
			// Path := c.Request.URL.Query().Get("Path")

			if ids == "" {
				c.JSON(404, "Empty Field")
			} else {

				client, ctx := mongoconnect.DBConnection()
				var createDB = client.Database(env.Data_Name).Collection("AddMathod")
				deletrezult, deleteerror := createDB.DeleteOne(ctx, bson.M{
					"_id": ids,
				})

				if deleteerror != nil {
					fmt.Printf("deleteerror: %v\n", deleteerror)
				}
				if deletrezult.DeletedCount == 1 {
					c.JSON(200, "succes")
					fmt.Printf("deletrezult: %v\n", deletrezult)
				} else {
					c.JSON(404, "error")
				}
			}
		}
	}
}