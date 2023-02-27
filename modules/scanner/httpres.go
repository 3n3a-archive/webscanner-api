package scanner

import (
	"io"
	"strings"
)

func (s *ScanClient) GetHTTPReponseInfo() (HttpResponseInfo, error) {
	// Get the supllied baseUrl's Headers
	resp, err := s.client.R().Get("")
	if err != nil || resp.IsErrorState() {
		return CustomOrDefaultError(
			"url couldn't be accessed",
			err,
			HttpResponseInfo{},
		)
	}

	copiedResCode := resp.StatusCode
	copiedHeaders := resp.Header.Clone()

	// Get Response Body
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return HttpResponseInfo{}, err
	}

	// Transform header keys to lowercase
	transformedHeaders := make(map[string][]string, 0)
	for name, values := range copiedHeaders {
		nameLower := strings.ToLower(name)
		transformedHeaders[nameLower] = values
	}

	// Get Response Headers
	return HttpResponseInfo{
		ResponseCode: copiedResCode,
		Headers:      transformedHeaders,
		TextBody:     string(body),
	}, nil
}