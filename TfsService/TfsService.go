package TfsService

type TfsProvider interface {
	GetRelations(ids TfsIds, linkType string) (TfsRelations, error)
	GetWorkItems(ids TfsIds, fields []string) (TfsWorkItems, error)
	UpdateWorkItem(id int, patch TfsPatchDocument) (TfsWorkItem, error)
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
		return workItems, nil
	}
	ids.AddTargets(relations)

	// Get child relations
	relations, err = s.Provider.GetRelations(ids, "System.LinkTypes.Hierarchy-Forward")
	if err != nil {
		return workItems, nil
	}
	ids.AddTargets(relations)

	// Get work items
	workItems, err = s.Provider.GetWorkItems(ids, fields)
	if err != nil {
		return workItems, nil
	}

	return workItems, nil
}

// Add work item comment
func (s *TfsService) AddWorkItemComment(id int, comment string) (TfsWorkItem, error) {
	var workItem TfsWorkItem

	if len(comment) > 1048576 {
		comment = comment[0 : 1048576-1]
	}

	operation := TfsPatchOperation{Op: "add", Path: "/fields/System.History", Value: comment}
	patchDocument := TfsPatchDocument{Operations: []TfsPatchOperation{operation}}

	workItem, err := s.Provider.UpdateWorkItem(id, patchDocument)
	if err != nil {
		return workItem, err
	}
	return workItem, nil
}
