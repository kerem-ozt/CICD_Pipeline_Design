package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kerem-ozt/GoodBlast_API/models"
	"github.com/kerem-ozt/GoodBlast_API/services"
)

// Leaderboard godoc
// @Summary      Init Leaderboard
// @Description  initing leaderboard
// @Tags         leaderboard
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /init [get]
// @Security     ApiKeyAuth
func EnsureLeaderboardInitialized(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	leaderboard, err := services.EnsureLeaderboardInitialized("global")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"leaderboard": leaderboard}
	response.SendResponse(c)
}

// Leaderboard godoc
// @Summary      Get Global Leaderboard
// @Description  get global leaderboard
// @Tags         leaderboard
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /getglobal [get]
// @Security     ApiKeyAuth
func GetGlobalLeaderboard(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	leaderboard, err := services.GetGlobalLeaderboard("global")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if len(leaderboard.Users) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "leaderboard is empty"})
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"leaderboard": leaderboard}
	response.SendResponse(c)
}

// Leaderboard godoc
// @Summary      Get Country Leaderboard
// @Description  get caountry leaderboard
// @Tags         leaderboard
// @Accept       json
// @Produce      json
// @Param        country  query
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /getcountry [get]
// @Security     ApiKeyAuth
func GetLeaderboardByCountry(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	country := c.Query("country")

	leaderboard, err := services.GetLeaderboardByCountry("global", country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if len(leaderboard.Users) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "leaderboard is empty"})
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"leaderboard": leaderboard}
	response.SendResponse(c)
}
