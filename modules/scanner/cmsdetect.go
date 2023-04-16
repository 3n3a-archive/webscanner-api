package scanner

import (
	"os"
	"regexp"
	"strings"

	file "github.com/3n3a/webscanner-api/modules/file"
	validate "github.com/3n3a/webscanner-api/modules/validation"

	// "github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
)

const (
	VERSION_REGEX = `(?m)([0-9.]{3,}|[0-9])([ \n"';,\-\t]{0,1})`
)



func (s *ScanClient) getVersionFromString(input string) string {
	// TODO: parse versions like '6.2-beta3-55420' --> kinda does
	regex := regexp.MustCompile(VERSION_REGEX)
	if matches := regex.FindStringSubmatch(input); cap(matches) > 0 {
		return matches[1]
	}
	return ""
}

func (s *ScanClient) getGeneratorTagTechnologies() ([]Technology, error) {
	GENERATOR_TECHNOLOGIES_YAML_FILE := os.Getenv("GENERATOR_TECHNOLOGIES_YAML_FILE")
	detectedTechnologies, err := file.ReadYAMLIntoStruct[[]Technology](GENERATOR_TECHNOLOGIES_YAML_FILE)
	if validate.IsErrorState(err) {
		return CustomOrDefaultError(
			"yaml technologies wasn't found",
			err,
			[]Technology{},
		)
	}


	// first just a singular example (one url, one factor)
	// TODO: use http response here
	resp, err := s.client.R().Get("")
	if validate.IsErrorState(err) || resp.IsErrorState() {
		return CustomOrDefaultError(
			"url couldn't be accessed",
			err,
			[]Technology{},
		)
	}

	defer resp.Body.Close()
	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return CustomOrDefaultError(
			"Error while accessing response",
			err,
			[]Technology{},
		)
	}

	// Find Generator Meta Tags
	foundGeneratorNames := make([]string, 0)
	metaTagsWithName := htmlquery.Find(doc, "//meta[@name]")
	for _, tag := range metaTagsWithName {
		name := strings.ToLower(htmlquery.SelectAttr(tag, "name"))

		if name == "generator" {
			content := htmlquery.SelectAttr(tag, "content")
			foundGeneratorNames = append(foundGeneratorNames, content)
		}
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

	return detectedTechnologies, nil
}

func (s *ScanClient) DetectTechnologies() (TechnologiesInfo, error) {
	// TODO: look at headers (x-powered-by, server, )
	var detectedTechnologies []Technology

	generatorTechnologies, err := s.getGeneratorTagTechnologies()
	if validate.IsErrorState(err) {
		return CustomOrDefaultError(
			"error while trying to determine technology of site",
			err,
			TechnologiesInfo{},
		)
	}

	detectedTechnologies = append(detectedTechnologies, generatorTechnologies...)

	// filtering to only the ones found
	technologiesInfo := TechnologiesInfo{}
	for _, tech := range detectedTechnologies {
		if tech.Score > 0 {
			technologiesInfo.Detected = append(technologiesInfo.Detected, tech)
		}
	}

	return technologiesInfo, nil
}