package scanner

import (
	RobotsTxtParser "github.com/3n3a/robotstxt-parser"
	SecurityTxtParser "github.com/3n3a/securitytxt-parser"
)

type ScanReport struct {
	SecurityTxt SecurityTxtParser.SecurityTxt
	RobotsTxt RobotsTxtParser.RobotsTxt
	HttpResponseInfo HttpResponseInfo
}

type HttpResponseInfo struct {
	Headers map[string][]string
	ResponseCode int
	TextBody string
}