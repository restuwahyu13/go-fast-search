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
		Limit               int64          `query:"limit" validate:"required,number,min=1,max=1000"`
		Page                int64          `query:"page" validate:"required,number,min=1"`
		FilterBy            string         `query:"filter_by" validate:"omitempty"`
		Filter              map[string]any `query:"filter" validate:"omitempty"`
		SortBy              string         `query:"sort_by" validate:"omitempty"`
		Sort                string         `query:"sort" validate:"omitempty,oneof=asc desc"`
		SearchBy            string         `query:"search_by" validate:"omitempty"`
		Search              string         `query:"search" validate:"omitempty"`
		HighlightAttributes string         `query:"highlight_attributes" validate:"omitempty"`
		MatchingStrategy    string         `query:"matching_strategy" validate:"omitempty,oneof=last all frequency"`
	}
)
