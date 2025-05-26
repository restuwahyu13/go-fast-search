package opt

type (
	MeiliSearchDocuments[T any] struct {
		Results    T      `json:"results,omitempty"`
		Hits       T      `json:"hits,omitempty"`
		Query      string `json:"query,omitempty"`
		Limit      int64  `json:"limit,omitempty"`
		Offset     int64  `json:"page,omitempty"`
		TotalPages int64  `json:"total_page,omitempty"`
		Total      int64  `json:"total,omitempty"`
	}
)
