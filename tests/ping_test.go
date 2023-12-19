package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kerem-ozt/GoodBlast_API/routes"
	"github.com/kerem-ozt/GoodBlast_API/services"
)

func TestUserService(t *testing.T) {
	services.LoadConfig()
	services.InitMongoDB()

	_, _ = services.EnsureLeaderboardInitialized("global")

	// if services.Config.UseRedis {
	// 	services.CheckRedisConnection()
	// }

	routes.InitGin()
	router := routes.New()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/v1/ping", nil)

	router.ServeHTTP(w, c.Request)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
