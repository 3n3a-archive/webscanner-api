package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	scanner "github.com/3n3a/webscanner-api/modules/scanner"
	validate "github.com/3n3a/webscanner-api/modules/validation"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func isErrorAddToList(err error, sR *scanner.ScanReport) {
	if err != nil {
		sR.Errors = append(sR.Errors, err.Error())
	}
}


type ScanRequestBody struct {
	BaseUrl string `json:"base_url"`
}


// TODO: check that only base-url was provided (e.g. host) or else parse from given url
func addScanRoutes(rg *gin.RouterGroup, customLogWriter *logrus.Logger) {
	scan := rg.Group("/scan")
	scan.POST("", func(c *gin.Context) {
		g := new(errgroup.Group)

		var reqBody ScanRequestBody
		if err := c.BindJSON(&reqBody); validate.IsErrorState(err) {
			validate.JsonError(err, http.StatusNotAcceptable, c)
			return
		}
		
		baseUrl := reqBody.BaseUrl;

		err := validate.ValidateUrl(baseUrl)
		if validate.IsErrorState(err) {
			validate.JsonError(err, http.StatusNotAcceptable, c)
			return
		}

		url, _ := url.Parse(baseUrl)
		cleanedBaseUrl := fmt.Sprintf("%s://%s", url.Scheme, url.Hostname())
		scanClient := scanner.ScanClient{}
		scanClient.Create("WebScanner/1.0", cleanedBaseUrl, customLogWriter)

		if !scanClient.HostExists(url) {
			c.JSON(http.StatusOK, scanner.ScanReport{
				Errors: []string{
					"Host does not exist",
				},
			})
			return
		}


		sR := scanner.ScanReport{}

		g.Go(func() error {
			start := startTimeCount()

			st, err := scanClient.GetSecurityTxt()
			sR.SecurityTxt = st
			isErrorAddToList(err, &sR)

			printElapsedTime(start, "Security.Txt", customLogWriter)

			return nil
		})

	
		g.Go(func() error {
			start := startTimeCount()

			rt, err := scanClient.GetRobotsTxt()
			sR.RobotsTxt = rt
			isErrorAddToList(err, &sR)
			
			printElapsedTime(start, "Robots.txt", customLogWriter)

			start2 := startTimeCount()

			sm, err := scanClient.GetSiteMaps()
			sR.SitemapIndexes = sm
			isErrorAddToList(err, &sR)

			printElapsedTime(start2, "Sitemaps", customLogWriter)

			return nil
		})

		g.Go(func() error {
			start := startTimeCount()

			hi, err := scanClient.GetHTTPReponseInfo()
			sR.HttpResponseInfo = hi
			isErrorAddToList(err, &sR)

			printElapsedTime(start, "Http Response Info", customLogWriter)

			return nil
		})


		g.Go(func() error {
			start := startTimeCount()

			tech, err := scanClient.DetectTechnologies()
			sR.Technologies = tech
			isErrorAddToList(err, &sR)

			printElapsedTime(start, "Technologies", customLogWriter)

			return nil
		})

		// NEW STUFF

		// END NEW

		if err := g.Wait(); err != nil {
			fmt.Println("Error while processing scan report")
		}

		c.JSON(http.StatusOK, sR)
	})
}

func startTimeCount() time.Time {
	return time.Now()
}

func printElapsedTime(t time.Time, name string, logger *logrus.Logger) {
	elapsed := time.Since(t)
	logger.Debugf("%s Took %s", name, elapsed)
}