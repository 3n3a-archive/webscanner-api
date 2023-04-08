package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func addPingRoutes(rg *gin.RouterGroup, logger *logrus.Logger) {
	ping := rg.Group("/ping")

	ping.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
}