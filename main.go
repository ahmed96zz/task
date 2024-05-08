package main

import "github.com/gin-gonic/gin"

func main() {
	initDB()
	router := gin.Default()
	router.POST("/api/users", creatUser)
	router.GET("/api/users", GetUsers)
	router.POST("/api/generateotp", createOTP)
	router.POST("/api/verifyotp", verifyOTP)
	router.Run(":8080")
}
