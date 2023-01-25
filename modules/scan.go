package scanner

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/3n3a/securitytxt-parser"
	"github.com/imroc/req/v3"
)

type ScanClient struct {
	baseUrl string

	client *req.Client
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

func (s *ScanClient) GetSecurityTxt() (SecurityTxtParser.SecurityTxt, error) {
	// Get the Security.Txt
	resp, err := s.client.R().
		Get("/.well-known/security.txt")
	if err != nil || resp.IsErrorState() {

		if err == nil {
			return SecurityTxtParser.SecurityTxt{}, errors.New("no security.txt found")
		} else {
			return SecurityTxtParser.SecurityTxt{}, err
		}
	}

	// Get Response Body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return SecurityTxtParser.SecurityTxt{}, err
	}

	// Parse .Txt
	inputTxt := string(body)
	st, err := SecurityTxtParser.ParseTxt(inputTxt)
	if err != nil {
		return SecurityTxtParser.SecurityTxt{}, err
	}

	return st, nil
}