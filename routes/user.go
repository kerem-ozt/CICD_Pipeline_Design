package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kerem-ozt/GoodBlast_API/controllers"
	"github.com/kerem-ozt/GoodBlast_API/middlewares/validators"
)

func UserRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	user := router.Group("/user", handlers...)
	{
		user.GET(
			"/whoami",
			controllers.WhoAmI,
		)

		user.GET(
			"/getall",
			controllers.GetAllUsers,
		)

		user.GET(
			"/getbyid",
			controllers.GetById,
		)

		user.DELETE(
			"/delete",
			controllers.DeleteUser,
		)

		user.POST(
			"/entertournament",
			controllers.EnterTournament,
		)

		user.POST(
			"/updateprogress",
			validators.UpdateProgress(),
			controllers.UpdateProgress,
		)
	}
}
