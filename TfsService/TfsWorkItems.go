package TfsService

type TfsFields struct {
	WorkItemType string `json:"System.WorkItemType"`
	State        string `json:"System.State"`
}

type TfsWorkItem struct {
	Id     int
	Fields TfsFields
}

type TfsWorkItems struct {
	Items []TfsWorkItem `json:"Value"`
}
