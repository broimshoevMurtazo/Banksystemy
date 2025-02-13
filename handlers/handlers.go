package handlers

import (
	"nank/app/controlers"
	"nank/app/createadmin"

	"github.com/gin-gonic/gin"
)

func Handlers() {
	createadmin.Createadmin()
	r := gin.Default()
	r.POST("/registration",controlers.Registration)
	r.POST("/Login",controlers.Login)
	r.POST("/AddMathod",controlers.AddMathod)
	r.POST("/updatepassword",controlers.UpdatePassword)
	r.POST("/Addmany",controlers.Income)
	r.DELETE("/DeleteMathod",controlers.DeleteMathod)
	r.POST("/UpdateMany",controlers.UpdateUserCash)
	r.GET("/verify",controlers.Verify)
	r.POST("/search",controlers.Search)
	r.POST("/newpass",controlers.ChangePass)


	
	r.Run(":2020")

}