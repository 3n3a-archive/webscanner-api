package scanner

import "io/ioutil"

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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return HttpResponseInfo{}, err
	}

	// Get Response Headers
	return HttpResponseInfo{
		ResponseCode: copiedResCode,
		Headers:      copiedHeaders,
		TextBody:     string(body),
	}, nil
}