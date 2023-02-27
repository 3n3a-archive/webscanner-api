package scanner

import (
	RobotsTxtParser "github.com/3n3a/robotstxt-parser"
	SecurityTxtParser "github.com/3n3a/securitytxt-parser"
)

type ResponseInterfaces interface {
	SecurityTxtParser.SecurityTxt | RobotsTxtParser.RobotsTxt | HttpResponseInfo | SitemapInfo | SitemapIndex | []SitemapIndex | TechnologiesInfo
}

type ScanReport struct {
	SecurityTxt SecurityTxtParser.SecurityTxt		`json:"security-txt"`
	RobotsTxt RobotsTxtParser.RobotsTxt				`json:"robots-txt"`
	HttpResponseInfo HttpResponseInfo				`json:"http-response-info"`
	SitemapIndexes []SitemapIndex					`json:"sitemapindexes"`
	Technologies TechnologiesInfo					`json:"technologies"`
	Errors []string 								`json:"errors"`
}

type HttpResponseInfo struct {
	Headers map[string][]string		`json:"headers"`
	ResponseCode int				`json:"response-code"`
	TextBody string 				`json:"-"`
}

type SitemapInfo struct {
	LocationUrl string				`json:"location-url"`
	Urls []string					`json:"urls"`
}

type SitemapIndex struct {
	Sitemaps []SitemapInfo			`json:"sitemaps"`
}

type TechnologiesInfo struct {
	Detected []Technology			`json:"detected"`
}

type Technology struct {
	DetectionString string 			`json:"-"`
	Name string						`json:"name"`
	Version string					`json:"version"`
	Score int						`json:"score"`
}