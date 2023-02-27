package scanner

import (
	"errors"
	"time"

	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
)

type ScanClient struct {
	baseUrl string
	logger *logrus.Logger

	client *req.Client

	sitemapUrls []string
}

func (s *ScanClient) Create(userAgent string, serverUrl string, logger *logrus.Logger) {
	s.logger = logger
	s.baseUrl = serverUrl

	s.client = req.C().
		SetBaseURL(s.baseUrl).
		SetCommonRetryCount(2).
		SetCommonRetryBackoffInterval(1*time.Millisecond, 100*time.Millisecond).
		SetCommonRetryFixedInterval(10*time.Millisecond).
		SetTimeout(5*time.Second).
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