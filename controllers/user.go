package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kerem-ozt/GoodBlast_API/models"
	db "github.com/kerem-ozt/GoodBlast_API/models/db"
	"github.com/kerem-ozt/GoodBlast_API/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateProgress godoc
// @Summary      Update Progress
// @Description  updates the progress of a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        req  body      models.ProgressRequest true "Progress Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /update [post]
// @Security     ApiKeyAuth
func UpdateProgress(c *gin.Context) {
	var requestBody models.ProgressRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userID := requestBody.UserID

	err := services.UpdateProgress(userID, requestBody.Score, requestBody.Coin)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.SendResponse(c)
}

// WhoAmI godoc
// @Summary      Who Am I
// @Description  returns the user information
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /whoami [get]
// @Security     ApiKeyAuth
func WhoAmI(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	token := c.GetHeader("Bearer-Token")
	tokenModel, err := services.VerifyToken(token, db.TokenTypeAccess)
	if err != nil {
		models.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := services.FindUserById(tokenModel.User)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"user": user}
	response.SendResponse(c)
}

// EnterTournament godoc
// @Summary      Enter Tournament
// @Description  enters the user to the tournament
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        tournamentID  query     string true  "Tournament ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /entertournament [get]
// @Security     ApiKeyAuth
func EnterTournament(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	token := c.GetHeader("Bearer-Token")
	tokenModel, err := services.VerifyToken(token, db.TokenTypeAccess)
	if err != nil {
		models.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := services.FindUserById(tokenModel.User)

	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	tournamentIDStr := c.Query("tournamentID")
	tournamentID, err := primitive.ObjectIDFromHex(tournamentIDStr)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	err = services.EnterTournament(user.ID, tournamentID)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.SendResponse(c)
}

// GetAllUsers godoc
// @Summary      Get All Users
// @Description  returns all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /getall [get]
// @Security     ApiKeyAuth
func GetAllUsers(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	users, err := services.GetAllUsers()
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"users": users}
	response.SendResponse(c)
}

// GetById godoc
// @Summary      Get User By Id
// @Description  returns user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  query     string true  "User ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /getbyid [get]
// @Security     ApiKeyAuth
func GetById(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	id := c.Query("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := services.FindUserById(objectID)
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"user": user}
	response.SendResponse(c)
}

// DeleteUser godoc
// @Summary      Delete User
// @Description  deletes user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  query     string true  "User ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /delete [get]
// @Security     ApiKeyAuth
func DeleteUser(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	id := c.Query("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = services.DeleteUser(objectID)
	if err != nil {
		models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.SendResponse(c)
}
