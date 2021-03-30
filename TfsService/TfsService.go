package TfsService

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Tfs sercvice
type TfsService struct {
	BaseUri            string
	АuthorizationToken string
	Client             *http.Client
}

// Do request
func (s *TfsService) DoRequest(method string, requestUrl string, body []byte) (string, error) {
	req, err := http.NewRequest(method, s.BaseUri+requestUrl+"?api-version=5.0", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(s.АuthorizationToken, s.АuthorizationToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Gets the results of the query given its WIQL
func (s *TfsService) QueryWiql(query string) (string, error) {
	body, err := json.Marshal(map[string]interface{}{"query": query})
	if err != nil {
		return "", err
	}

	res, err := s.DoRequest("POST", "_apis/wit/wiql", body)
	if err != nil {
		return "", err
	}
	return res, nil
}

// Get work items batch
func (s *TfsService) GetWorkItemsBatch(ids TfsIds, fields []string) (string, error) {
	body, err := json.Marshal(map[string]interface{}{"ids": ids.Ids, "fields": fields, "errorPolicy": "Omit"})
	if err != nil {
		return "", err
	}

	res, err := s.DoRequest("POST", "_apis/wit/workitemsbatch", body)
	if err != nil {
		return "", err
	}
	return res, nil
}
