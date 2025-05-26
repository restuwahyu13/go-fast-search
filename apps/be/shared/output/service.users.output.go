package opt

import entitie "github.com/restuwahyu13/go-fast-search/domain/entities"

type (
	UsersSearch struct {
		Results []entitie.UsersDocument `json:"results"`
		Limit   int64                   `json:"limit"`
		Offset  int64                   `json:"offset"`
		Total   int64                   `json:"total"`
	}
)
