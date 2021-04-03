package TfsService

type TfsPatchOperation struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

type TfsPatchDocument struct {
	Operations []TfsPatchOperation
}
