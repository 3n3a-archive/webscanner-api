package scanner

import (
	"io/ioutil"
	"reflect"

	SecurityTxtParser "github.com/3n3a/securitytxt-parser"
)

func (s *ScanClient) GetSecurityTxt() (SecurityTxtParser.SecurityTxt, error) {
	// Get the Security.Txt
	resp, err := s.client.R().
		Get("/.well-known/security.txt")
	if err != nil || resp.IsErrorState() {
		return CustomOrDefaultError(
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
		return CustomOrDefaultError(
			"no security.txt found",
			err,
			SecurityTxtParser.SecurityTxt{},
		)
	}

	return st, nil
}
