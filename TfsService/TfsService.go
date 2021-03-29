package TfsService

import (
	"BugNetSyncService/Tools"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const ApiVersion = "api-version=5.0"

// Tfs sercvice
type TfsService struct {
	BaseUri            string
	АuthorizationToken string
	Client             *http.Client
}

// Gets the results of the query given its WIQL
func (s *TfsService) QueryWiql(query string) (string, error) {
	var result string

	body, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest("POST", s.BaseUri+"_apis/wit/wiql?"+ApiVersion, bytes.NewBuffer(body))
	if err != nil {
		return result, err
	}
	req.SetBasicAuth(s.АuthorizationToken, s.АuthorizationToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	result = string(bytes)

	return result, nil
}

// Get work items relations
func (s *TfsService) GetWorkItemsRelations(tfsIds []int, linkType string) (string, error) {
	query := "SELECT System.Id FROM WorkItemLinks WHERE Source.System.Id in(" + Tools.JoinToString(tfsIds, ",") + ") and System.Links.LinkType='" + linkType + "'"

	result, err := s.QueryWiql(query)
	if err != nil {
		return result, err
	}
	return result, nil
}
