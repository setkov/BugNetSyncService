package TfsService

import (
	"BugNetSyncService/Common"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Tfs provider
type tfsProvider struct {
	BaseUri            string
	АuthorizationToken string
	Client             *http.Client
}

// New Tfs provider
func NewTfsProvider(baseUri string, authorizationToken string) *tfsProvider {
	return &tfsProvider{
		BaseUri:            baseUri,
		АuthorizationToken: authorizationToken,
		Client:             &http.Client{},
	}
}

// Do request
func (p *tfsProvider) doRequest(method string, requestUrl string, body []byte) (string, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s?api-version=5.0", p.BaseUri, requestUrl), bytes.NewBuffer(body))
	if err != nil {
		return "", Common.NewError("New request. " + err.Error())
	}
	req.SetBasicAuth(p.АuthorizationToken, p.АuthorizationToken)

	contentType := "application/json"
	if method == "PATCH" {
		contentType += "-patch+json"
	}
	req.Header.Add("Content-Type", contentType)

	resp, err := p.Client.Do(req)
	if err != nil {
		return "", Common.NewError("Do request. " + err.Error())
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", Common.NewError("Read request. " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		var response map[string]interface{}
		_ = json.Unmarshal(bytes, &response)
		return "", Common.NewError(fmt.Sprintf("Request status %v. %v", resp.Status, response["message"]))
	}
	return string(bytes), nil
}

// Gets the results of the query given its WIQL
func (p *tfsProvider) queryWiql(query string) (string, error) {
	body, err := json.Marshal(map[string]interface{}{"query": query})
	if err != nil {
		return "", Common.NewError("Prepare for query wiql. " + err.Error())
	}

	res, err := p.doRequest("POST", "_apis/wit/wiql", body)
	if err != nil {
		return "", Common.NewError("Query wiql. " + err.Error())
	}
	return res, nil
}

// Get work items batch
func (p *tfsProvider) getWorkItemsBatch(ids TfsIds, fields []string) (string, error) {
	body, err := json.Marshal(map[string]interface{}{"ids": ids.Ids, "fields": fields, "errorPolicy": "Omit"})
	if err != nil {
		return "", Common.NewError("Prepare for get work items batch. " + err.Error())
	}

	res, err := p.doRequest("POST", "_apis/wit/workitemsbatch", body)
	if err != nil {
		return "", Common.NewError("Get work items batch. " + err.Error())
	}
	return res, nil
}

// Get relations
func (p *tfsProvider) GetRelations(ids TfsIds, linkType string) (TfsRelations, error) {
	var relations TfsRelations

	query := "SELECT System.Id FROM WorkItemLinks WHERE Source.System.Id in(" + ids.JoinToString(",") + ") and System.Links.LinkType='" + linkType + "'"
	res, err := p.queryWiql(query)
	if err != nil {
		return relations, Common.NewError("Get relations. " + err.Error())
	}

	err = json.Unmarshal([]byte(res), &relations)
	if err != nil {
		return relations, Common.NewError("Parse relations. " + err.Error())
	}
	return relations, nil
}

// Get work items
func (p *tfsProvider) GetWorkItems(ids TfsIds, fields []string) (TfsWorkItems, error) {
	var workItems TfsWorkItems

	res, err := p.getWorkItemsBatch(ids, fields)
	if err != nil {
		return workItems, Common.NewError("Get work items. " + err.Error())
	}

	err = json.Unmarshal([]byte(res), &workItems)
	if err != nil {
		return workItems, Common.NewError("Parse work items. " + err.Error())
	}
	return workItems, nil
}

// Update work item
func (p *tfsProvider) UpdateWorkItem(id int, patch TfsPatchDocument) (TfsWorkItem, error) {
	var workItem TfsWorkItem

	body, err := json.Marshal(patch.Operations)
	if err != nil {
		return workItem, Common.NewError("Prepare for update work item. " + err.Error())
	}

	res, err := p.doRequest("PATCH", fmt.Sprintf("_apis/wit/workitems/%d", id), body)
	if err != nil {
		return workItem, Common.NewError("Update work item. " + err.Error())
	}

	err = json.Unmarshal([]byte(res), &workItem)
	if err != nil {
		return workItem, Common.NewError("Parse work item. " + err.Error())
	}
	return workItem, nil
}
