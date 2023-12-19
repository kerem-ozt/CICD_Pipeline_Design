package validators

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kerem-ozt/GoodBlast_API/models"
)

func CreateNewTournament() gin.HandlerFunc {
	return func(c *gin.Context) {

		var tournamentRequest models.TournamentRequest
		_ = c.ShouldBindBodyWith(&tournamentRequest, binding.JSON)

		if err := tournamentRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}
