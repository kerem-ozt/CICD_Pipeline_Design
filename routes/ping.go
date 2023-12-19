package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kerem-ozt/GoodBlast_API/controllers"
)

func PingRoute(router *gin.RouterGroup) {
	auth := router.Group("/ping")
	{
		auth.GET(
			"",
			controllers.Ping,
		)
	}
}
