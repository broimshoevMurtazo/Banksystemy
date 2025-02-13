package controlers

import (
	"context"
	"fmt"
	env "nank/app/Env"
	hashedpasswod "nank/app/HashedPassword"
	htmldata "nank/app/HtmlData"
	mongoconnect "nank/app/MongoConnect"
	emptyfieldcheker "nank/app/emptyfiledcheker"
	"nank/app/halpers"
	returnjwt "nank/app/jwt"
	"nank/app/sendmail"
	"nank/app/structs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Registration(c *gin.Context) {
	var userdata structs.Registration
	if err := c.ShouldBindJSON(&userdata); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input data"})
		return
	}

	// Проверяем на пустые поля
	emptyfield, err := emptyfieldcheker.EmptyField(userdata, "Id", "Permission", "SecrateCode", "Cash", "verify")
	if emptyfield {
		c.JSON(400, gin.H{"error": err})
		return
	}

	client, ctx := mongoconnect.DBConnection()
	defer client.Disconnect(ctx)

	connect := client.Database(env.Data_Name).Collection("Users")

	findrezult := connect.FindOne(ctx, bson.M{
		"email": userdata.Email,
	})
	var Dbdata structs.Registration
	findrezult.Decode(&Dbdata)

	if Dbdata.Email == "" && Dbdata.Phone == "" {
		// Генерируем код и ID
		Code := halpers.RandomSixDigit()
		verify := false
		ID := primitive.NewObjectID().Hex()
		userdata.SecrateCode = Code

		// Хэшируем пароль
		Hashpass, _ := hashedpasswod.HashPassword(userdata.Password)

		// Вставляем пользователя в коллекцию Users
		insertResult, insertError := connect.InsertOne(ctx, bson.M{
			"_id":         ID,
			"name":        userdata.Name,
			"surname":     userdata.Surname,
			"password":    Hashpass,
			"email":       userdata.Email,
			"phone":       userdata.Phone,
			"permission":  "user",
			"secrateCode": userdata.SecrateCode,
			"verify":      verify,
			"cash":        0,
		})

		if insertError != nil {
			fmt.Printf("Insert error: %v\n", insertError)
			c.JSON(500, gin.H{"error": "Failed to register user"})
			return
		}

		fmt.Printf("Insert result: %v\n", insertResult)

		// Отправляем письмо пользователю
		htmldata.HtmlData(userdata)
		sendmail.SendGomail("controllers/htmls/output.html", "Your secret code", userdata.Email)

		// Создаем коллекцию code
		connect2 := client.Database(env.Data_Name).Collection("code")
		id2 := primitive.NewObjectID().Hex()
		_, insertError2 := connect2.InsertOne(ctx, bson.M{
			"_id":         id2,
			"secrateCode": userdata.SecrateCode,
			"email":       userdata.Email,
		})

		if insertError2 != nil {
			fmt.Printf("Insert error (code): %v\n", insertError2)
			c.JSON(500, gin.H{"error": "Failed to generate verification code"})
			return
		}

		// Удаляем код через 60 секунд
		time.AfterFunc(60*time.Second, func() {
			fmt.Printf("Attempting to delete code with ID: %s\n", id2)
			Delete(id2)
		})

		c.JSON(200, gin.H{"message": "Registration successful"})
	} else {
		c.JSON(399, "email or phone alrady exsisted!!!")
	}
}
func Delete(ids string) {
	client, ctx := mongoconnect.DBConnection()
	defer client.Disconnect(ctx)

	connect2 := client.Database(env.Data_Name).Collection("code")
	deleteResult, deleteError := connect2.DeleteOne(ctx, bson.M{
		"_id": ids,
	})

	if deleteError != nil {
		fmt.Printf("Delete error: %v\n", deleteError)
		return
	}

	if deleteResult.DeletedCount == 1 {
		fmt.Println("Code successfully deleted")
	} else {
		fmt.Println("Code not found")
	}
}

func Login(c *gin.Context) {
	var logintempt structs.Registration

	c.ShouldBindJSON(&logintempt)

	Emptyfield, err := emptyfieldcheker.EmptyField(logintempt, "Name", "Surname", "Id", "Phone", "Permission", "SecrateCode", "Cash", "verify")
	if Emptyfield {
		c.JSON(404, err)
	} else {
		client, ctx := mongoconnect.DBConnection()
		DBconnect := client.Database(env.Data_Name).Collection("Users")
		result := DBconnect.FindOne(ctx, bson.M{
			"email": logintempt.Email,
		})
		var Userdata structs.Registration
		result.Decode(&Userdata)
		if Userdata.Email != "" {
			IsValidPass := hashedpasswod.CompareHashPasswords(Userdata.Password, logintempt.Password)
			fmt.Printf("IsValidPass: %v\n", IsValidPass)
			key := returnjwt.GenerateToken(Userdata.Phone, Userdata.Permission, Userdata.Id)
			fmt.Printf("key: %v\n", key)
			if IsValidPass {
				http.SetCookie(c.Writer, &http.Cookie{
					Name:     env.Data_Name,
					Value:    key,
					Expires:  time.Now().Add(60 * time.Hour),
					Domain:   "",
					Path:     "/",
					Secure:   false,
					HttpOnly: false,
					SameSite: http.SameSiteLaxMode,
				})
				c.JSON(200, "success")
			} else {
				c.JSON(404, "Not valid pass")
			}
		} else {
			c.JSON(350, "User not found")
		}
	}
}
func Income(c *gin.Context) {
	var cookidata, cookieerror = c.Request.Cookie(env.Data_Name)
	if cookieerror != nil {
		fmt.Printf("cookieerror: %v\n", cookieerror)
		c.JSON(404, "error Not Cookie found")
		fmt.Printf("cookidata: %v\n", cookidata)
	} else {
		SecretKeyData, _ := returnjwt.Validate(cookidata.Value)
		if SecretKeyData.Permission == "user" || SecretKeyData.Permission == "admin" {
			var CashTempt structs.CashStruct
			c.ShouldBindJSON(&CashTempt)
			fmt.Printf("Cash: %v\n", CashTempt)

			Emptyfiled, err := emptyfieldcheker.EmptyField(CashTempt, "Id", "User_ID", "Phone", "Sender_id", "Reciver_id")
			if Emptyfiled {
				c.JSON(404, err)
			} else {
				client, ctx := mongoconnect.DBConnection()
				connect := client.Database(env.Data_Name).Collection("Users")
				findresult2 := connect.FindOne(ctx, bson.M{
					"phone": CashTempt.SenderNumber,
				})
				var sender structs.CashStruct

				findresult2.Decode(&sender)

				if sender.Phone == "" {
					c.JSON(404, "Error phone not founded")
				} else {
					if sender.Cash < CashTempt.Cash {
						c.JSON(300, "there is no anoth many")
					} else {
						findresult3 := connect.FindOne(ctx, bson.M{
							"phone": CashTempt.ReciverNumber,
						})
						var recivertempt structs.CashStruct

						findresult3.Decode(&recivertempt)
						_, err3 := connect.UpdateOne(ctx,
							bson.M{
								"phone": CashTempt.SenderNumber,
							},
							bson.D{
								{Key: "$set", Value: bson.M{
									"cash": sender.Cash - CashTempt.Cash,
								},
								},
							})
						if err3 != nil {
							fmt.Printf("err3: %v\n", err3)
						} else {
							connect2 := client.Database(env.Data_Name).Collection("Income")
							_, einserterror := connect2.InsertOne(ctx, bson.M{
								"_id":           primitive.NewObjectID().Hex(),
								"ReciverNumber": CashTempt.ReciverNumber,
								"Reciver_id":    recivertempt.Id,
								"cash":          CashTempt.Cash,
								"SenderNumber":  CashTempt.SenderNumber,
								"Sender_id":     sender.Id,
							})
							if einserterror != nil {
								fmt.Println("error")
							}
						}
						findresult := connect.FindOne(ctx, bson.M{
							"phone": CashTempt.ReciverNumber,
						})
						var dbdata structs.Registration
						findresult.Decode(&dbdata)
						fmt.Println(dbdata)
						if dbdata.Email != "" {
							new_cash := dbdata.Cash + CashTempt.Cash

							_, err2 := connect.UpdateOne(ctx,
								bson.M{
									"phone": CashTempt.ReciverNumber,
								},
								bson.D{
									{Key: "$set", Value: bson.M{
										"cash": new_cash,
									},
									},
								})
							connect3 := client.Database(env.Data_Name).Collection("cheks")
							_, einserterror2 := connect3.InsertOne(ctx, bson.M{
								"_id":           primitive.NewObjectID().Hex(),
								"ReciverNumber": CashTempt.ReciverNumber,
								"sender_id":     sender.Id,
								"reciver_id":    recivertempt.Id,
								"cash":          CashTempt.Cash,
								"SenderNumber":  CashTempt.SenderNumber,
							})
							if einserterror2 != nil {
								fmt.Printf("einserterror: %v\n", einserterror2)
							} else {
								c.JSON(200, "succes")

							}
							if err2 != nil {
								fmt.Printf("err2: %v\n", err2)
							}

						}
					}

				}

			}

		}
	}

}
func UpdateUserCash(c *gin.Context) {
	var cookidata, cookieerror = c.Request.Cookie(env.Data_Name)
	if cookieerror != nil {
		fmt.Printf("cookieerror: %v\n", cookieerror)
		c.JSON(404, "error Not Cookie found")
		fmt.Printf("cookidata: %v\n", cookidata)
	} else {
		SecretKeyData, _ := returnjwt.Validate(cookidata.Value)
		if SecretKeyData.Permission == "user" || SecretKeyData.Permission == "admin" {
			var AddMany structs.AddCash
			c.ShouldBindJSON(&AddMany)

			Emptyfiled, err := emptyfieldcheker.EmptyField(AddMany, "Id")
			if Emptyfiled {
				c.JSON(404, err)
			} else {
				client, ctx := mongoconnect.DBConnection()
				connect := client.Database(env.Data_Name).Collection("Users")
				findrezult := connect.FindOne(context.TODO(), bson.M{
					"phone": AddMany.Phone,
				})
				var UpdateCash structs.AddCash

				findrezult.Decode(&UpdateCash)

				if UpdateCash.Phone != "" {
					_, err2 := connect.UpdateOne(ctx,
						bson.M{
							"phone": UpdateCash.Phone,
						},
						bson.D{
							{Key: "$set", Value: bson.M{
								"cash": UpdateCash.Cash + AddMany.Cash,
							},
							},
						})
					if err2 != nil {
						fmt.Println(err2)
						fmt.Println("Update error first")
					} else {
						connect2 := client.Database(env.Data_Name).Collection("addmanny")
						_, inserterror := connect2.InsertOne(ctx, bson.M{
							"_id":    primitive.NewObjectID().Hex(),
							"phone":  UpdateCash.Phone,
							"income": AddMany.Cash,
						})
						if inserterror != nil {
							fmt.Println(inserterror)
						} else {
							c.JSON(200, "succes")
						}
					}
				} else {
					c.JSON(390, "phone not founded")
				}
			}
		}
	}
}

//
//

func AddMathod(c *gin.Context) {
	Cockiedata, cokieerror := c.Request.Cookie(env.Data_Name)
	if cokieerror != nil {
		c.JSON(404, "error not cokie found")
	} else {
		Secretkey, Isvalid := returnjwt.Validate(Cockiedata.Value)
		fmt.Println(Secretkey.Permission, "aasas")
		if Secretkey.Permission != "Admin" && Isvalid {
			c.JSON(404, "eroor")
		} else {
			var addmathod structs.AddMathod
			c.ShouldBindJSON(&addmathod)
			Emptyfield, err := emptyfieldcheker.EmptyField(addmathod, "Id")
			if Emptyfield {
				c.JSON(404, err)
			} else {
				client, ctx := mongoconnect.DBConnection()
				Id := primitive.NewObjectID().Hex()
				Dbconnect := client.Database(env.Data_Name).Collection("AddMathod")
				insertrezult, inserterror := Dbconnect.InsertOne(ctx, bson.M{
					"_id":        Id,
					"Name":       addmathod.Name,
					"Descrition": addmathod.Description,
					"Logo":       addmathod.Logo,
				})
				if inserterror != nil {
					fmt.Printf("inserterror: %v\n", inserterror)
				} else {
					fmt.Printf("insertrezult: %v\n", insertrezult)
					c.JSON(200, "succes")
				}
			}
		}
	}
}

func Search(c *gin.Context) {
	var searchingData structs.Registration

	c.ShouldBindJSON(&searchingData)
	if searchingData.Phone == "" && searchingData.Password == "" {
		c.JSON(400, "error")
	}
	client, ctx := mongoconnect.DBConnection()
	connect := client.Database(env.Data_Name).Collection("Users")
	finddata := connect.FindOne(ctx, bson.M{
		"phone": searchingData.Phone,
	})
	var newdata structs.Registration
	finddata.Decode(newdata)
	hashedpasswod.CompareHashPasswords(searchingData.Password,newdata.Password)
	if newdata.Password != searchingData.Password {
		c.JSON(404,"error the passwrod for these user phone is incorect!!!")
	} else {
		collection := client.Database(env.Data_Name).Collection("cheks")
		filter := bson.D{
			{"$match", bson.D{
				{"SenderNumber", bson.D{
					{"$regex", searchingData.Phone},
					{"$options", "i"},
				}},
			}},
		}
		cursor, err := collection.Aggregate(ctx, mongo.Pipeline{filter})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
			return
		}
		defer cursor.Close(ctx)
		var results []bson.M
		if err := cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read query results"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}
