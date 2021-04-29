package TfsService

import (
	"BugNetSyncService/Common"
)

type TfsProvider interface {
	GetRelations(ids TfsIds, linkType string) (TfsRelations, error)
	GetWorkItems(ids TfsIds, fields []string) (TfsWorkItems, error)
	UpdateWorkItem(id int, patch TfsPatchDocument) (TfsWorkItem, error)
	CreateAttachment(fileName string, content []byte) (AttachmentReference, error)
}

// Tfs sercvice
type TfsService struct {
	Provider TfsProvider
}

// New tfs service
func NewTfsService(p TfsProvider) *TfsService {
	return &TfsService{Provider: p}
}

// Get work items with related
func (s *TfsService) GetWorkItemsRelated(ids TfsIds, fields []string) (TfsWorkItems, error) {
	var workItems TfsWorkItems

	// Get related relations
	relations, err := s.Provider.GetRelations(ids, "System.LinkTypes.Related")
	if err != nil {
		return workItems, Common.NewError("Get related. " + err.Error())
	}
	ids.AddTargets(relations)

	// Get child relations
	relations, err = s.Provider.GetRelations(ids, "System.LinkTypes.Hierarchy-Forward")
	if err != nil {
		return workItems, Common.NewError("Get child. " + err.Error())
	}
	ids.AddTargets(relations)

	// Get work items
	workItems, err = s.Provider.GetWorkItems(ids, fields)
	if err != nil {
		return workItems, Common.NewError("Get work items related. " + err.Error())
	}

	return workItems, nil
}

// Add work item comment
func (s *TfsService) AddWorkItemComment(id int, comment string) (TfsWorkItem, error) {
	var workItem TfsWorkItem

	if len(comment) > 1048576 {
		comment = comment[0 : 1048576-1]
	}

	operation := TfsPatchOperation{
		Op:    "add",
		Path:  "/fields/System.History",
		Value: comment,
	}
	patchDocument := TfsPatchDocument{Operations: []TfsPatchOperation{operation}}

	workItem, err := s.Provider.UpdateWorkItem(id, patchDocument)
	if err != nil {
		return workItem, Common.NewError("Add work item comment. " + err.Error())
	}
	return workItem, nil
}

// Add work item attachment
func (s *TfsService) AddWorkItemAttachment(id int, fileName string, content []byte) (TfsWorkItem, error) {
	var workItem TfsWorkItem

	ref, err := s.Provider.CreateAttachment(fileName, content)
	if err != nil {
		return workItem, Common.NewError("Create attachment. " + err.Error())
	}

	atachmentValue := AttachmentValue{
		Rel: "AttachedFile",
		Url: ref.Url,
		Attributes: AttachmentAttributes{
			Comment: "Attached with BugNet Integration Tools",
		},
	}
	operation := TfsPatchOperation{
		Op:    "add",
		Path:  "/relations/-",
		Value: atachmentValue,
	}
	patchDocument := TfsPatchDocument{Operations: []TfsPatchOperation{operation}}

	workItem, err = s.Provider.UpdateWorkItem(id, patchDocument)
	if err != nil {
		return workItem, Common.NewError("Add work item comment. " + err.Error())
	}
	return workItem, nil
}
