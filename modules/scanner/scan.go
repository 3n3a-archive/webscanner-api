package scanner

import (
	"errors"
	"io/ioutil"
	"reflect"
	"time"

	"github.com/3n3a/robotstxt-parser"
	"github.com/3n3a/securitytxt-parser"
	"github.com/imroc/req/v3"
)

type ResponseInterfaces interface {
	SecurityTxtParser.SecurityTxt | RobotsTxtParser.RobotsTxt | HttpResponseInfo
}

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

	return rt, nil
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