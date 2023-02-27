package scanner

import (
	"errors"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
)

type ScanClient struct {
	baseUrl string
	logger *logrus.Logger

	client *req.Client

	remoteAddr net.IP

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

func (s *ScanClient) HostExists(url *url.URL) bool {
	address := url.Host + ":80"
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)

// 	fmt.Printf("Connection established between %s and localhost with time out of %d seconds.\n", address, 5)
//    fmt.Printf("Remote Address : %s \n", conn.RemoteAddr().String())
//    fmt.Printf("Local Address : %s \n", conn.LocalAddr().String())

	if err != nil {
		return false
	}

	s.remoteAddr = net.ParseIP(strings.Split(conn.RemoteAddr().String(), ":")[0])

	return err == nil
}

func CustomOrDefaultError[S ResponseInterfaces](message string, defaultError error, emptyStruct S) (S, error) {
	if defaultError == nil {
		return emptyStruct, errors.New(message)
	} else {
		return emptyStruct, defaultError
	}
}