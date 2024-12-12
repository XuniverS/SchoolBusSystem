package main

import (
	"awesomeProject/back"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(cors.Default())
	router.Static("", "./front")

	back.InitDatabase()

	back.RegisterUserModule(router)
	back.RegisterSetupRoutes(router)
	back.RegisterProfileModule(router)

	router.Run(":8000")

}
