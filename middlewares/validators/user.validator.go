package validators

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kerem-ozt/GoodBlast_API/models"
)

func UpdateProgress() gin.HandlerFunc {
	return func(c *gin.Context) {

		var progressRequest models.ProgressRequest
		_ = c.ShouldBindBodyWith(&progressRequest, binding.JSON)

		if err := progressRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}
