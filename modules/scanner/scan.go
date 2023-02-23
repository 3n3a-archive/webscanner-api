package scanner

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	"github.com/3n3a/robotstxt-parser"
	"github.com/3n3a/securitytxt-parser"
	"github.com/imroc/req/v3"
	"github.com/oxffaa/gopher-parse-sitemap"
)

type ResponseInterfaces interface {
	SecurityTxtParser.SecurityTxt | RobotsTxtParser.RobotsTxt | HttpResponseInfo | SitemapInfo | SitemapIndex | []SitemapIndex
}

type ScanClient struct {
	baseUrl string

	client *req.Client

	sitemapUrls []string
}

func (s *ScanClient) Create(userAgent string, serverUrl string) {
	s.baseUrl = serverUrl

	s.client = req.C().
		SetBaseURL(s.baseUrl).
		SetCommonRetryCount(2).
		SetCommonRetryBackoffInterval(1 * time.Millisecond, 100 * time.Millisecond).
    	SetCommonRetryFixedInterval(10 * time.Millisecond).
		SetTimeout(5 * time.Second).
		SetCommonHeader("Accept", "application/json").
		SetUserAgent(userAgent)
		
		// DevMode().
		// EnableDumpEachRequest()
}

func customOrDefaultError[S ResponseInterfaces](message string, defaultError error, emptyStruct S) (S, error) {
	if defaultError == nil {
		return emptyStruct, errors.New(message)
	} else {
		return emptyStruct, defaultError
	}
}

func (s *ScanClient) GetSecurityTxt() (SecurityTxtParser.SecurityTxt, error) {
	// Get the Security.Txt
	resp, err := s.client.R().
		Get("/.well-known/security.txt")
	if err != nil || resp.IsErrorState() {
		return customOrDefaultError(
			"no security.txt found",
			err,
			SecurityTxtParser.SecurityTxt{},
		)
	}

	// Get Response Body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return SecurityTxtParser.SecurityTxt{}, err
	}

	// Parse .Txt
	inputTxt := string(body)
	st, err := SecurityTxtParser.ParseTxt(inputTxt)
	if err != nil || reflect.DeepEqual(st, SecurityTxtParser.SecurityTxt{}) {
		return customOrDefaultError(
			"no security.txt found",
			err,
			SecurityTxtParser.SecurityTxt{},
		)
	}

	return st, nil
}

func (s *ScanClient) GetRobotsTxt() (RobotsTxtParser.RobotsTxt, error) {
	// Get the Robots.txt
	resp, err := s.client.R().
		Get("/robots.txt")
	if err != nil || resp.IsErrorState() {
		return customOrDefaultError(
			"no robots.txt found",
			err,
			RobotsTxtParser.RobotsTxt{},
		)
	}

	// Get Response Body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RobotsTxtParser.RobotsTxt{}, err
	}

	// Parse .Txt
	inputTxt := string(body)
	rt, err := RobotsTxtParser.ParseTxt(inputTxt)
	if err != nil || reflect.DeepEqual(rt, RobotsTxtParser.RobotsTxt{}) {
		return customOrDefaultError(
			"no robots.txt found",
			err,
			RobotsTxtParser.RobotsTxt{},
		)
	}


	// Detect if has Sitemaps
	if len(rt.Sitemaps) > 0 {
		s.sitemapUrls = append(s.sitemapUrls, rt.Sitemaps...)
	}

	return rt, nil
}

func (s *ScanClient) isSitemapIndex(body io.Reader) bool {
	bodyStart := make([]byte, 256)
	_, err := body.Read(bodyStart)
	if err != nil {
		return false
	}

	bodyStartString := string(bodyStart)
	return strings.Contains(bodyStartString, "<sitemapindex")

}

func (s *ScanClient) getSitemapUrlsByUrl(url string) []string {
	var urls []string
	err := sitemap.ParseFromSite(url, func(e sitemap.Entry) error {
		urls = append(urls, e.GetLocation())
		return nil
	})
	if err != nil {
		return urls
	}

	return urls
}

func (s *ScanClient) getSitemap(bodyBuffer io.Reader, originUrl string) SitemapInfo {
	var currentSitemap SitemapInfo
	currentSitemap.LocationUrl = originUrl

	err := sitemap.Parse(bodyBuffer, func(e sitemap.Entry) error {
		currentSitemap.Urls = append(currentSitemap.Urls, e.GetLocation())
		return nil
	})
	if err != nil {
		return currentSitemap
	}

	return currentSitemap
}

func (s *ScanClient) getSitemapIndex(bodyBuffer io.Reader) SitemapIndex {
	var currentIndex SitemapIndex
	var sitemapsCounter int
	err := sitemap.ParseIndex(bodyBuffer, func(e sitemap.IndexEntry) error {
		// Only get top 10 sitemaps from index
		if sitemapsCounter <= 10 {
			currentIndex.Sitemaps = append(currentIndex.Sitemaps, SitemapInfo{
				LocationUrl: e.GetLocation(),
				Urls: s.getSitemapUrlsByUrl(e.GetLocation()),
			})
			sitemapsCounter++
		}
		return nil
	})
	if err != nil {
		return currentIndex
	}

	return currentIndex
}

func (s *ScanClient) sitemapExists(sitemapUrl string) bool {
	resp, err := req.C().R().Get(sitemapUrl)
	if err != nil || resp.IsErrorState() {
		fmt.Println(err, resp.IsErrorState())
		return false
	}

	return true
}
 
func (s *ScanClient) GetSiteMaps() ([]SitemapIndex, error) {
	// Get the file
	if len(s.sitemapUrls) == 0 {
		sitemapUrlString := fmt.Sprintf("%s/%s", s.baseUrl, "sitemap.xml")
		if !s.sitemapExists(sitemapUrlString) {
			return make([]SitemapIndex, 0), nil
		}

		// else continue on
		s.sitemapUrls = append(s.sitemapUrls, sitemapUrlString)
	}

	var sitemapIndexes []SitemapIndex

	for _, sitemapUrl := range s.sitemapUrls {
		resp, err := req.C().R().Get(sitemapUrl)
		if err != nil || resp.IsErrorState() {
			return customOrDefaultError(
				"sitemap couldn't be accessed",
				err,
				make([]SitemapIndex, 0),
			)
		}
			
		// Read Body into Memory
		// This might be dangerous
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return make([]SitemapIndex, 0), err
		}
		
		bodyBuffer := bytes.NewBuffer(body)
		bodyBuffer2 := bytes.NewBuffer(body)
	
		if s.isSitemapIndex(bodyBuffer2) {
			sitemapIndexes = append(sitemapIndexes, s.getSitemapIndex(bodyBuffer))
		} else {
			// todo: maybe eventually check if is acuallty a sitemap
			sitemaps := make([]SitemapInfo, 0)
			sitemaps = append(sitemaps, s.getSitemap(bodyBuffer, sitemapUrl), )
			sitemapIndexes = append(sitemapIndexes, SitemapIndex{
				Sitemaps: sitemaps,
			})
		}

	}
	

	return sitemapIndexes, nil
}

func (s *ScanClient) GetHTTPReponseInfo() (HttpResponseInfo, error) {
	// Get the supllied baseUrl's Headers
	resp, err := s.client.R().Get("")
	if err != nil || resp.IsErrorState() {
		return customOrDefaultError(
			"url couldn't be accessed",
			err,
			HttpResponseInfo{},
		)
	}

	copiedResCode := resp.StatusCode
	copiedHeaders := resp.Header.Clone()
	
	// Get Response Body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return HttpResponseInfo{}, err
	}

	// Get Response Headers
	return HttpResponseInfo{
		ResponseCode: copiedResCode,
		Headers: copiedHeaders,
		TextBody: string(body),
	}, nil
}