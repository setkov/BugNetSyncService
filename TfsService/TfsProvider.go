package TfsService

import "encoding/json"

// Tfs provider
type TfsProvider struct {
	TfsService *TfsService
}

// New Tfs privider
func NewTfsProvider(s *TfsService) TfsProvider {
	return TfsProvider{
		TfsService: s,
	}
}

// Get relations
func (p *TfsProvider) GetRelations(ids TfsIds, linkType string) (TfsRelations, error) {
	var relations TfsRelations

	query := "SELECT System.Id FROM WorkItemLinks WHERE Source.System.Id in(" + ids.JoinToString(",") + ") and System.Links.LinkType='" + linkType + "'"
	res, err := p.TfsService.QueryWiql(query)
	if err != nil {
		return relations, err
	}

	err = json.Unmarshal([]byte(res), &relations)
	if err != nil {
		return relations, err
	}
	return relations, nil
}

// Get work items
func (p *TfsProvider) GetWorkItems(ids TfsIds, fields []string) (TfsWorkItems, error) {
	var workItems TfsWorkItems

	res, err := p.TfsService.GetWorkItemsBatch(ids, fields)
	if err != nil {
		return workItems, err
	}

	err = json.Unmarshal([]byte(res), &workItems)
	if err != nil {
		return workItems, err
	}
	return workItems, nil
}

// Get work items with related
func (p *TfsProvider) GetWorkItemsRelated(ids TfsIds, fields []string) (TfsWorkItems, error) {
	var workItems TfsWorkItems

	// Get related relations
	relations, err := p.GetRelations(ids, "System.LinkTypes.Related")
	if err != nil {
		return workItems, nil
	}
	ids.AddTargets(relations)

	// Get child relations
	relations, err = p.GetRelations(ids, "System.LinkTypes.Hierarchy-Forward")
	if err != nil {
		return workItems, nil
	}
	ids.AddTargets(relations)

	// Get work items
	workItems, err = p.GetWorkItems(ids, fields)
	if err != nil {
		return workItems, nil
	}

	return workItems, nil
}
