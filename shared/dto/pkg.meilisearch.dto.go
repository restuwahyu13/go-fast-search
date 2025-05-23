package dto

type MeiliSearchDocuments[T any] struct {
	ID     any    `json:"id"`
	Doc    string `json:"doc"`
	Data   T      `json:"data"`
	IsBulk bool   `json:"is_bulk"`
	Action string `json:"action"`
}
