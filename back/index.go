package back

import "github.com/gin-gonic/gin"

func RegisterIndexModule(router *gin.Engine) {
	indexRouter := router.Group("/index")
	{
		indexRouter.POST("/queryAll")
		indexRouter.POST("/book")
		indexRouter.POST("/payed")
		indexRouter.POST("/unbook")
		indexRouter.POST("/queryBooked")
		indexRouter.POST("/queryFinished")
	}
}

func queryAll() {
	
}
