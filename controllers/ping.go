package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kerem-ozt/GoodBlast_API/models"
)

// Ping godoc
// @Summary      Ping
// @Description  check server
// @Tags         ping
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Router       /ping [get]
func Ping(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "pong",
	}

	response.SendResponse(c)
}
