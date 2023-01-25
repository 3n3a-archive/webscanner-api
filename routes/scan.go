package routes

import (
	"net/http"
	"net/url"

	scanner "github.com/3n3a/webscanner-api/modules"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func addScanRoutes(rg *gin.RouterGroup) {
	scan := rg.Group("/scan")

	
	scan.GET("/", func(c *gin.Context) {
		baseUrl := c.Query("base_url")
		parsedUrl, err := url.Parse(baseUrl)
		allowedSchemes := []string{"http", "https"}
		if baseUrl == "" || 
			err != nil || 
			!slices.Contains(allowedSchemes, parsedUrl.Scheme) {

			c.JSON(http.StatusNotAcceptable, gin.H{
				"message": "Invalid Base Url. Please enter a url",
			})
			return
		}

		scanClient := scanner.ScanClient{}
		scanClient.Create("WebScanner/1.0", baseUrl)
		st, err := scanClient.GetSecurityTxt()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, st)
	})
}