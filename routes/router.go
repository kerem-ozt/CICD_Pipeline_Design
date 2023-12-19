package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kerem-ozt/GoodBlast_API/docs"
	"github.com/kerem-ozt/GoodBlast_API/middlewares"
	"github.com/kerem-ozt/GoodBlast_API/models"
	"github.com/kerem-ozt/GoodBlast_API/services"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New() *gin.Engine {
	r := gin.New()
	initRoute(r)

	r.Use(gin.LoggerWithWriter(middlewares.LogWriter()))
	r.Use(gin.CustomRecovery(middlewares.AppRecovery()))
	r.Use(middlewares.CORSMiddleware())

	v1 := r.Group("/v1")
	{
		PingRoute(v1)
		AuthRoute(v1)
		UserRoute(v1, middlewares.JWTMiddleware())
		TournamentRoute(v1, middlewares.JWTMiddleware())
		leaderboardRouter(v1, middlewares.JWTMiddleware())
	}

	docs.SwaggerInfo.BasePath = v1.BasePath() // adds /v1 to swagger base path

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

func initRoute(r *gin.Engine) {
	_ = r.SetTrustedProxies(nil)
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true

	r.NoRoute(func(c *gin.Context) {
		models.SendErrorResponse(c, http.StatusNotFound, c.Request.RequestURI+" not found")
	})

	r.NoMethod(func(c *gin.Context) {
		models.SendErrorResponse(c, http.StatusMethodNotAllowed, c.Request.Method+" is not allowed here")
	})
}

func InitGin() {
	gin.DisableConsoleColor()
	gin.SetMode(services.Config.Mode)
}
