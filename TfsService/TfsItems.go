package TfsService

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
	WorkItemType TfsWorkItemType
	State        TfsWorkItemState
}
