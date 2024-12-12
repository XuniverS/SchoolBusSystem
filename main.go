package main

import (
	"awesomeProject/back"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// 注册中间件，例如跨域中间件
	router.Use(cors.Default())
	back.RegisterloginModule(router)

	router.Run(":8000")

}
