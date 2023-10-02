package main

import (
	"jwt/controller"
	"jwt/middleware"
	"jwt/server"
)



func main(){
	s := server.NewServer()

	s.R.POST("/signup", controller.SignUp)
	s.R.POST("/login", controller.Login)
	s.R.GET("/validate", middleware.RequireAUth,controller.Validate)

	s.R.Run(":8080")


		

}