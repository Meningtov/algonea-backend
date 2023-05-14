package api

import (
	"net/http"

	"github.com/Meningtov/algonea-backend/handler"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func registerRoutes(r *gin.RouterGroup) {
	r.GET("/api/health", handler.HealthCheck)
	r.GET("/api/account/:address/send-asa", handler.SendAsa)
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	app = gin.New()
	r := app.Group("/")
	registerRoutes(r)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
