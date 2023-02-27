package scanner

import (
	"io/ioutil"
	"reflect"

	RobotsTxtParser "github.com/3n3a/robotstxt-parser"
)

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