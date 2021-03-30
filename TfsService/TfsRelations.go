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
