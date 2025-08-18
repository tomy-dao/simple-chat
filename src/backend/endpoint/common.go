package endpoint

type Response[T any] struct {
	Data  *T      `json:"data"`
	Error string `json:"error"`
}