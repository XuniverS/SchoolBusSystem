package main

import (
	"awesomeProject/backend"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	backend.InitDatabase()

	backend.RegisterUserModule(router)
	backend.RegisterSetupRoutes(router)
	backend.RegisterProfileModule(router)
	backend.RegisterIndexModule(router)

	router.Static("/htmls", "./front/htmls")

	router.Run(":8000")

}
