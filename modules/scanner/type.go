package scanner

import (
	RobotsTxtParser "github.com/3n3a/robotstxt-parser"
	SecurityTxtParser "github.com/3n3a/securitytxt-parser"
)

type ResponseInterfaces interface {
	SecurityTxtParser.SecurityTxt | RobotsTxtParser.RobotsTxt | HttpResponseInfo | SitemapInfo | SitemapIndex | []SitemapIndex | TechnologiesInfo
}

type ScanReport struct {
	SecurityTxt SecurityTxtParser.SecurityTxt		`json:"securitytxt"`
	RobotsTxt RobotsTxtParser.RobotsTxt				`json:"robotstxt"`
	HttpResponseInfo HttpResponseInfo				`json:"httpresponseinfo"`
	SitemapIndexes []SitemapIndex					`json:"sitemapindexes"`
	Technologies TechnologiesInfo					`json:"technologies"`
	Errors []string 								`json:"errors"`
}

type HttpResponseInfo struct {
	Headers map[string][]string		`json:"headers"`
	ResponseCode int				`json:"responsecode"`
	TextBody string 				`json:"-"`
}

type SitemapInfo struct {
	LocationUrl string				`json:"locationurl"`
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