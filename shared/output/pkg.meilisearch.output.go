package opt

type (
	MeiliSearchDocuments[T any] struct {
		Results            T      `json:"results,omitempty"`
		Hits               T      `json:"hits,omitempty"`
		EstimatedTotalHits int64  `json:"estimatedTotalHits,omitempty"`
		Offset             int64  `json:"offset,omitempty"`
		Limit              int64  `json:"limit,omitempty"`
		ProcessingTimeMs   int64  `json:"processingTimeMs,omitempty"`
		Query              string `json:"query,omitempty"`
		FacetDistribution  any    `json:"facetDistribution,omitempty"`
		TotalHits          int64  `json:"totalHits,omitempty"`
		HitsPerPage        int64  `json:"hitsPerPage,omitempty"`
		Page               int64  `json:"page,omitempty"`
		TotalPages         int64  `json:"totalPages,omitempty"`
		FacetStats         any    `json:"facetStats,omitempty"`
		IndexUID           string `json:"indexUid,omitempty"`
		Total              int64  `json:"total,omitempty"`
	}
)
