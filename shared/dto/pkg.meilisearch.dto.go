package dto

type (
	MeiliSearchDocuments[T any] struct {
		ID     any    `json:"id"`
		Doc    string `json:"doc"`
		Data   T      `json:"data"`
		IsBulk bool   `json:"is_bulk"`
		Action string `json:"action"`
	}

	MeiliSearchDocumentsQuery struct {
		Limit            int64          `query:"limit" validate:"required,number,min=1,max=1000"`
		Page             int64          `query:"page" validate:"required,number,min=1"`
		Filter           map[string]any `query:"filter" validate:"omitempty"`
		Sort             string         `query:"sort" validate:"omitempty,oneof=asc desc"`
		Search           string         `query:"search" validate:"omitempty"`
		MatchingStrategy string         `query:"matching_strategy" validate:"omitempty,oneof=last all frequency"`
	}
)
