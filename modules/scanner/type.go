package scanner

import (
	RobotsTxtParser "github.com/3n3a/robotstxt-parser"
	SecurityTxtParser "github.com/3n3a/securitytxt-parser"
)

type ResponseInterfaces interface {
	SecurityTxtParser.SecurityTxt | RobotsTxtParser.RobotsTxt | HttpResponseInfo | SitemapInfo | SitemapIndex | []SitemapIndex | TechnologiesInfo
}

type ScanReport struct {
	SecurityTxt SecurityTxtParser.SecurityTxt
	RobotsTxt RobotsTxtParser.RobotsTxt
	HttpResponseInfo HttpResponseInfo
	SitemapIndexes []SitemapIndex
	Technologies TechnologiesInfo
	Errors []string
}

type HttpResponseInfo struct {
	Headers map[string][]string
	ResponseCode int
	TextBody string
}

type SitemapInfo struct {
	LocationUrl string
	Urls []string
}

type SitemapIndex struct {
	Sitemaps []SitemapInfo
}

type TechnologiesInfo struct {
	Detected []Technology
}

type Technology struct {
	DetectionString string `json:"-"`
	Name string
	Version string
	Score int
}