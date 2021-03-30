package TfsService

import "encoding/json"

type TfsWorkItemType string

const (
	Bug         TfsWorkItemType = "Bug"
	Feature     TfsWorkItemType = "Feature"
	Issue       TfsWorkItemType = "Issue"
	Requirement TfsWorkItemType = "Requirement"
)

type TfsWorkItemState string

const (
	Active   TfsWorkItemState = "Active"
	Closed   TfsWorkItemState = "Closed"
	Proposed TfsWorkItemState = "Proposed"
	Resolved TfsWorkItemState = "Resolved"
)

type TfsFields struct {
	WorkItemType string `json:"System.WorkItemType"`
	State        string `json:"System.State"`
	Title        string `json:"System.Title"`
}

type TfsWorkItem struct {
	Id     int
	Fields TfsFields
}

type TfsWorkItems struct {
	Items []TfsWorkItem `json:"Value"`
}

// Load work items
func (i *TfsWorkItems) Load(s *TfsService, ids TfsIds, fields []string) error {
	res, err := s.GetWorkItemsBatch(ids, fields)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(res), &i)
	if err != nil {
		return err
	}
	return nil
}
