package scanner

import (
	"regexp"
	"strings"

	// "github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
)

const (
	VERSION_REGEX = `(?m)([0-9.]{3,})([ \n"';,\-\t]{0,1})`
)

func (s *ScanClient) getVersionFromString(input string) string {
	// TODO: parse versions like '6.2-beta3-55420' --> kinda does
	regex := regexp.MustCompile(VERSION_REGEX)
	if matches := regex.FindStringSubmatch(input); cap(matches) > 0 {
		return matches[1]
	}
	return ""
}

func (s *ScanClient) DetectTechnologies() (TechnologiesInfo, error) {
	// TODO: look at headers (x-powered-by, server, )

	// right now this array of techs is only designed for generator tag
	detectedTechnologies := []Technology {
		{
			DetectionString: "wordpress",
			Name: "WordPress",
			Version: "",
			Score: 0,
		},
		{
			DetectionString: "woocommerce",
			Name: "WooCommerce",
			Version: "",
			Score: 0,
		},
		{
			DetectionString: "wpml",
			Name: "WPML",
			Version: "",
			Score: 0,
		},
		{
			DetectionString: "elementor",
			Name: "Elementor",
			Version: "",
			Score: 0,
		},
		{
			DetectionString: "powered by wpbakery page builder",
			Name: "WPBakery Page Builder",
			Version: "",
			Score: 0,
		},
		{
			DetectionString: "all in one seo",
			Name: "All in One SEO (AIOSEO)",
			Version: "",
			Score: 0,
		},
	}


	// first just a singular example (one url, one factor)
	resp, err := s.client.R().Get("")
	if err != nil || resp.IsErrorState() {
		return CustomOrDefaultError(
			"url couldn't be accessed",
			err,
			TechnologiesInfo{},
		)
	}

	defer resp.Body.Close()
	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return CustomOrDefaultError(
			"Error while accessing response",
			err,
			TechnologiesInfo{},
		)
	}

	// Find Generator Meta Tags
	foundGeneratorNames := make([]string, 0)
	metaTagsWithName := htmlquery.Find(doc, "//meta[@name='generator']")
	for _, tag := range metaTagsWithName {
		// name := htmlquery.SelectAttr(tag, "name")
		content := htmlquery.SelectAttr(tag, "content")
		foundGeneratorNames = append(foundGeneratorNames, content)
	}


	// Check names against map and add scores
	for _, nameString := range foundGeneratorNames {
		for index, tech := range detectedTechnologies {
			nameLower := strings.ToLower(nameString)
			if strings.HasPrefix(nameLower, tech.DetectionString) {
				detectedTechnologies[index].Score += 1
				detectedTechnologies[index].Version = s.getVersionFromString(nameString)
			}
		}
	}

	// filtering to only the ones found
	technologiesInfo := TechnologiesInfo{}
	for _, tech := range detectedTechnologies {
		if tech.Score > 0 {
			technologiesInfo.Detected = append(technologiesInfo.Detected, tech)
		}
	}

	return technologiesInfo, nil
}