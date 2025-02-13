package controlers

import (
	"fmt"
	env "nank/app/Env"
	hashedpasswod "nank/app/HashedPassword"
	htmldata "nank/app/HtmlData"
	mongoconnect "nank/app/MongoConnect"
	emptyfieldcheker "nank/app/emptyfiledcheker"
	"nank/app/genertaepasswprd"
	returnjwt "nank/app/jwt"
	"nank/app/sendmail"
	"nank/app/structs"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdatePassword(c *gin.Context) {
	var update structs.Registration
	c.ShouldBindJSON(&update)
	emptyfield, err := emptyfieldcheker.EmptyField(update, "Name", "Surname", "Phone", "Pasword", "SecrateCode", "Permission", "Id", "Password", "Cash")
	if emptyfield {
		c.JSON(404, err)
	} else {
		client, ctx := mongoconnect.DBConnection()

		connect := client.Database(env.Data_Name).Collection("Users")

		findresult := connect.FindOne(ctx, bson.M{
			"email": update.Email,
		})
		var secondone structs.Registration
		findresult.Decode(&secondone)
		pass, _ := genertaepasswprd.GeneratePassword(8)

		if secondone.Email != "" {
			secondone.Password = pass
			hash, _ := hashedpasswod.HashPassword(secondone.Password)
			_, err2 := connect.UpdateOne(ctx,
				bson.M{
					"email": secondone.Email,
				},
				bson.D{
					{Key: "$set", Value: bson.M{
						"password": hash,
					},
					},
				})
			if err2 != nil {
				fmt.Printf("err2: %v\n", err2)
			} else {
				htmldata.HtmlData2(secondone)
				err := sendmail.SendGomail("./controlers/htmls/passwordoutput.html", "your password ", secondone.Email)
				if err != nil {
					fmt.Printf("err: %v\n", err)
				} else {
					c.JSON(200, "succes")
				}
			}
		} else {
			c.JSON(404, "email not found")
		}
	}
}
func ChangePass(c *gin.Context) {
	var cookidata, cookieerror = c.Request.Cookie(env.Data_Name)
	if cookieerror != nil {
		fmt.Printf("cookieerror: %v\n", cookieerror)
		c.JSON(404, "error Not Cookie found")
		fmt.Printf("cookidata: %v\n", cookidata)
	} else {
		SecretKeyData, _ := returnjwt.Validate(cookidata.Value)
		if SecretKeyData.Permission == "user" {
			var updatePass structs.UpdatePass
			c.ShouldBindJSON(&updatePass)
			emptyfiled, emptyerror := emptyfieldcheker.EmptyField(updatePass)
			if emptyfiled {
				c.JSON(404, emptyerror)
			} else {
				client, ctx := mongoconnect.DBConnection()
				connect := client.Database(env.Data_Name).Collection("Users")
				findresult := connect.FindOne(ctx, bson.M{
					"email": updatePass.Email,
				})
				var newdata structs.UpdatePass
				findresult.Decode(&newdata)
				if newdata.Email != "" {
					isvallidpass := hashedpasswod.CompareHashPasswords(newdata.Password, updatePass.Password)

					if isvallidpass {
						hash,_:=hashedpasswod.HashPassword(updatePass.Newpassword)
						connect.UpdateOne(ctx, bson.M{
							"email": updatePass.Email,
						},
							bson.D{
								{
									Key: "$set", Value: bson.M{
										"password": hash,
									},
								},
							},
						)
						c.JSON(200, "succes")
					} else {
						c.JSON(404, "incorrect password")
					}
				}else {
					c.JSON(404,"email not found")
				}

			}

		}
	}
}
