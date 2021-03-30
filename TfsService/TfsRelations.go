package TfsService

import "encoding/json"

type WorkItemId struct {
	Id int
}

type TfsRelation struct {
	Source WorkItemId
	Target WorkItemId
}

type TfsRelations struct {
	WorkItemRelations []TfsRelation
}

// Load relations
func (r *TfsRelations) Load(s *TfsService, tfsIds TfsIds, linkType string) error {
	query := "SELECT System.Id FROM WorkItemLinks WHERE Source.System.Id in(" + tfsIds.JoinToString(",") + ") and System.Links.LinkType='" + linkType + "'"
	res, err := s.QueryWiql(query)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(res), &r)
	if err != nil {
		return err
	}
	return nil
}
