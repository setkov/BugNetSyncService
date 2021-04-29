package TfsService

type AttachmentAttributes struct {
	Comment string `json:"comment"`
}

type AttachmentValue struct {
	Rel        string               `json:"rel"`
	Url        string               `json:"url"`
	Attributes AttachmentAttributes `json:"attributes"`
}

type AttachmentReference struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}
