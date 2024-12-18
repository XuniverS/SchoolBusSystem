package main

import (
	"awesomeProject/backend"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func init() {
	backend.InitDatabase()
}

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	pprof.Register(router, "debug/pprof")

	backend.RegisterUserModule(router)
	backend.RegisterSetupRoutes(router)
	backend.RegisterProfileModule(router)
	backend.RegisterIndexModule(router)

	router.Static("/htmls", "./front/htmls")
	router.Static("/img", "./front/img")

	router.Run(":8000")

}
