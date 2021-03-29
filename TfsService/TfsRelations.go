package TfsService

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

type TfsIds struct {
	Ids []int
}

// test TfsIds.Ids contains element
func (i *TfsIds) Contains(element int) bool {
	for _, id := range i.Ids {
		if id == element {
			return true
		}
	}
	return false
}

// add distinct target ids
func (i *TfsIds) AddTargets(relations TfsRelations) {
	for _, relation := range relations.WorkItemRelations {
		if !i.Contains(relation.Target.Id) {
			i.Ids = append(i.Ids, relation.Target.Id)
		}
	}
}
