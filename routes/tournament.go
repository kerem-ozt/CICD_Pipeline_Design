package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kerem-ozt/GoodBlast_API/controllers"
	"github.com/kerem-ozt/GoodBlast_API/middlewares/validators"
)

func TournamentRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	tournaments := router.Group("/tournament", handlers...)
	{
		tournaments.POST(
			"/create",
			validators.CreateNewTournament(),
			controllers.CreateNewTournament,
		)

		tournaments.POST(
			"/creategroup",
			controllers.CreateTournamentGroups,
		)

		tournaments.GET(
			"/getall",
			controllers.GetTournaments,
		)

		tournaments.GET(
			"/getbyid",
			controllers.GetTournamentById,
		)

		tournaments.POST(
			"/progress",
			controllers.ProgressTournament,
		)

		tournaments.GET(
			"/gettournamentresults",
			controllers.GetTournamentWinnersFromCache,
		)
	}
}
