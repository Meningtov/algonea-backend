package api

import (
	"github.com/Meningtov/algonea_backend/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func registerRoutes(r *gin.RouterGroup) {
	r.GET("/api/health", handler.HealthCheck)
	r.GET("/api/send-asa", handler.SendAsa)
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
