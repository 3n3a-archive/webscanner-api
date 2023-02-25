package scanner

import (
	"errors"
	"time"

	RobotsTxtParser "github.com/3n3a/robotstxt-parser"
	SecurityTxtParser "github.com/3n3a/securitytxt-parser"
	"github.com/imroc/req/v3"
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
		SetCommonRetryBackoffInterval(1*time.Millisecond, 100*time.Millisecond).
		SetCommonRetryFixedInterval(10*time.Millisecond).
		SetTimeout(5*time.Second).
		SetCommonHeader("Accept", "application/json").
		SetUserAgent(userAgent)

	// DevMode().
	// EnableDumpEachRequest()
}

func CustomOrDefaultError[S ResponseInterfaces](message string, defaultError error, emptyStruct S) (S, error) {
	if defaultError == nil {
		return emptyStruct, errors.New(message)
	} else {
		return emptyStruct, defaultError
	}
}