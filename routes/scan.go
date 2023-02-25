package routes

import (
	"net/http"

	scanner "github.com/3n3a/webscanner-api/modules/scanner"
	validate "github.com/3n3a/webscanner-api/modules/validation"
	"github.com/gin-gonic/gin"
)

func isErrorAddToList(err error, sR *scanner.ScanReport) {
	if err != nil {
		sR.Errors = append(sR.Errors, err.Error())
	}
}

func addScanRoutes(rg *gin.RouterGroup) {
	scan := rg.Group("/scan")
	scan.GET("", func(c *gin.Context) {
		baseUrl := c.Query("base_url")
		err := validate.ValidateUrl(baseUrl)
		if validate.IsErrorState(err) {
			validate.JsonError(err, http.StatusNotAcceptable, c)
			return
		}

		scanClient := scanner.ScanClient{}
		scanClient.Create("WebScanner/1.0", baseUrl)

		sR := scanner.ScanReport{}

		
		st, err := scanClient.GetSecurityTxt()
		sR.SecurityTxt = st
		isErrorAddToList(err, &sR)

	
	
		rt, err := scanClient.GetRobotsTxt()
		sR.RobotsTxt = rt
		isErrorAddToList(err, &sR)


		sm, err := scanClient.GetSiteMaps()
		sR.SitemapIndexes = sm
		isErrorAddToList(err, &sR)


	
		hi, err := scanClient.GetHTTPReponseInfo()
		sR.HttpResponseInfo = hi
		isErrorAddToList(err, &sR)


		c.JSON(http.StatusOK, sR)
	})

	
	scan.GET("/securitytxt", func(c *gin.Context) {
		baseUrl := c.Query("base_url")
		err := validate.ValidateUrl(baseUrl)
		if validate.IsErrorState(err) {
			validate.JsonError(err, http.StatusNotAcceptable, c)
			return
		}

		scanClient := scanner.ScanClient{}
		scanClient.Create("WebScanner/1.0", baseUrl)
		st, err := scanClient.GetSecurityTxt()

		if validate.IsErrorState(err) {
			validate.JsonError(err, http.StatusBadRequest, c)
			return
		}
		c.JSON(http.StatusOK, st)
	})

	scan.GET("/robotstxt", func(c *gin.Context) {
		baseUrl := c.Query("base_url")
		err := validate.ValidateUrl(baseUrl)
		if validate.IsErrorState(err) {
			validate.JsonError(err, http.StatusNotAcceptable, c)
			return
		}

		scanClient := scanner.ScanClient{}
		scanClient.Create("WebScanner/1.0", baseUrl)
		st, err := scanClient.GetRobotsTxt()

		if validate.IsErrorState(err) {
			validate.JsonError(err, http.StatusBadRequest, c)
			return
		}
		c.JSON(http.StatusOK, st)
	})
}