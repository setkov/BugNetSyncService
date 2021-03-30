package TfsService

import (
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
	body, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.BaseUri+"_apis/wit/wiql?"+ApiVersion, bytes.NewBuffer(body))
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

// Get work items relations
func (s *TfsService) GetWorkItemsRelations(tfsIds TfsIds, linkType string) (TfsRelations, error) {
	var rel = TfsRelations{}

	query := "SELECT System.Id FROM WorkItemLinks WHERE Source.System.Id in(" + tfsIds.JoinToString(",") + ") and System.Links.LinkType='" + linkType + "'"
	res, err := s.QueryWiql(query)
	if err != nil {
		return rel, err
	}

	err = json.Unmarshal([]byte(res), &rel)
	if err != nil {
		return rel, err
	}

	return rel, nil
}
