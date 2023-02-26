package scanner

import (
	"io/ioutil"
	"reflect"

	RobotsTxtParser "github.com/3n3a/robotstxt-parser"
)

// TODO: fix all the false positives
/*
This came out:

####################################################
"UserAgentRules": [
{
	"UserAgents": [
		"ia_archiver",
		"*",
		"WebReaper",
		"WebCopier",
		"Offline Explorer",
		"HTTrack",
		"Microsoft.URL.Control",
		"EmailCollector",
		"penthesilea"
	],
	"Allow": [
		"/"
	],
	"Disallow": [
		"/"
	]
}
####################################################

while this was put in (it's done the wronge pairing up):
####################################################
# ===================================
# Webseite: http://dfjdflkjdlf.ch/
# ===================================

Sitemap: http://www.djfldjfljlkdf.ch/sitemap.xml

User-agent: ia_archiver
Disallow: /

User-Agent: *
Allow: /

# ===================================
# SchlieÃŸe folgende Spider komplett aus:
# ===================================

User-agent: WebReaper
User-agent: WebCopier
User-agent: Offline Explorer
User-agent: HTTrack
User-agent: Microsoft.URL.Control
User-agent: EmailCollector
User-agent: penthesilea
####################################################3

*/
func (s *ScanClient) GetRobotsTxt() (RobotsTxtParser.RobotsTxt, error) {
	// Get the Robots.txt
	resp, err := s.client.R().
		Get("/robots.txt")
	if err != nil || resp.IsErrorState() {
		return CustomOrDefaultError(
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
		return CustomOrDefaultError(
			"no robots.txt found",
			err,
			RobotsTxtParser.RobotsTxt{},
		)
	}

	// Detect if has Sitemaps
	if cap(rt.Sitemaps) > 0 {
		s.sitemapUrls = append(s.sitemapUrls, rt.Sitemaps...)
	}

	return rt, nil
}