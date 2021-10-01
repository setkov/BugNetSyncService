package TfsService

import (
	"BugNetSyncService/Common"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiVersion = "api-version=5.0"

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
func (p *tfsProvider) doRequest(method string, requestUrl string, contentType string, body []byte) (string, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", p.BaseUri, requestUrl), bytes.NewBuffer(body))
	if err != nil {
		return "", Common.NewError("New request. " + err.Error())
	}
	req.SetBasicAuth(p.АuthorizationToken, p.АuthorizationToken)
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
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

	requestUrl := fmt.Sprintf("_apis/wit/wiql?%s", apiVersion)
	res, err := p.doRequest("POST", requestUrl, "application/json", body)
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

	requestUrl := fmt.Sprintf("_apis/wit/workitemsbatch?%s", apiVersion)
	res, err := p.doRequest("POST", requestUrl, "application/json", body)
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
	var workItems, workItemsPart TfsWorkItems

	//  get work items by 200 items
	var size int = 200
	var idsPage int = 0
	var idsPart = ids.TakePart(size, idsPage*size)
	for len(idsPart.Ids) > 0 {
		res, err := p.getWorkItemsBatch(idsPart, fields)
		if err != nil {
			return workItems, Common.NewError("Get work items. " + err.Error())
		}

		err = json.Unmarshal([]byte(res), &workItemsPart)
		if err != nil {
			return workItems, Common.NewError("Parse work items. " + err.Error())
		}
		workItems.Items = append(workItems.Items, workItemsPart.Items...)

		idsPage += 1
		idsPart = ids.TakePart(size, idsPage*size)
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

	requestUrl := fmt.Sprintf("_apis/wit/workitems/%d?%s", id, apiVersion)
	res, err := p.doRequest("PATCH", requestUrl, "application/json-patch+json", body)
	if err != nil {
		return workItem, Common.NewError("Update work item. " + err.Error())
	}

	err = json.Unmarshal([]byte(res), &workItem)
	if err != nil {
		return workItem, Common.NewError("Parse work item. " + err.Error())
	}
	return workItem, nil
}

// Create attachment
func (p *tfsProvider) CreateAttachment(fileName string, content []byte) (AttachmentReference, error) {
	var attachmentReference AttachmentReference

	// body, err := json.Marshal(content)
	// if err != nil {
	// 	return attachmentReference, Common.NewError("Prepare for create attachment. " + err.Error())
	// }
	//log.Print(body)

	requestUrl := fmt.Sprintf("_apis/wit/attachments?fileName=%s&%s", url.QueryEscape(fileName), apiVersion)
	res, err := p.doRequest("POST", requestUrl, "application/octet-stream", content)
	if err != nil {
		return attachmentReference, Common.NewError("Create attachment. " + err.Error())
	}

	err = json.Unmarshal([]byte(res), &attachmentReference)
	if err != nil {
		return attachmentReference, Common.NewError("Parse attachment reference. " + err.Error())
	}
	return attachmentReference, nil
}
