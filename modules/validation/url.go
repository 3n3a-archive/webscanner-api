package validate

import (
	"errors"
	"net/url"

	"golang.org/x/exp/slices"
)

func ValidateUrl(urlToValidate string) error {
	parsedUrl, err := url.Parse(urlToValidate)
	allowedSchemes := []string{"http", "https"}
	if urlToValidate == "" || 
		err != nil || 
		!slices.Contains(allowedSchemes, parsedUrl.Scheme) {

		return errors.New("invalid url. please enter a valid url")
	}
	return nil
}