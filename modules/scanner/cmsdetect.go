package scanner

import (
	"regexp"
	"strings"

	// "github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
)

func (s *ScanClient) getVersionFromString(input string) string {
	// TODO: parse versions like '6.2-beta3-55420'
	regex := regexp.MustCompile(`(?m)([0-9.]{5,})`)
	return regex.FindString(input)
}

func (s *ScanClient) DetectTechnologies() (TechnologiesInfo, error) {
	detectedTechnologies := []Technology {
		{
			Name: "WordPress",
			Version: "",
			Score: 0,
		},
		{
			Name: "WooCommerce",
			Version: "",
			Score: 0,
		},
		{
			Name: "WPML",
			Version: "",
			Score: 0,
		},
		{
			Name: "Elementor",
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
			techLower := strings.ToLower(tech.Name)
			if strings.Contains(nameLower, techLower) {
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