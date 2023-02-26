package routes

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	logrus "github.com/sirupsen/logrus"
)

var customLogWriter = logrus.New()

var router = gin.New()	

// Run will start the server
func Run() {
	// Middleware
	customLogWriter.Formatter = &logrus.JSONFormatter{}

	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.BestCompression))
	router.Use(gin.LoggerWithWriter(customLogWriter.Writer()))

	// Routes
	getRoutes()
	router.Run(":5000")
}

// getRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func getRoutes() {
	v1 := router.Group("/v1")
	addScanRoutes(v1)
	addPingRoutes(v1)
}