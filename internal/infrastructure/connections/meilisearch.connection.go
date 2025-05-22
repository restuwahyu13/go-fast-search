package con

import (
	"crypto/tls"
	"net/http"

	"github.com/meilisearch/meilisearch-go"

	"github.com/restuwahyu13/go-fast-search/shared/dto"
)

func NewMeiliSearch(req dto.Request[dto.Environtment]) meilisearch.ServiceManager {
	maxRetries := 15

	return meilisearch.New(req.Config.MEILISEARCH.URL,
		meilisearch.WithAPIKey(req.Config.MEILISEARCH.KEY),
		meilisearch.WithCustomClientWithTLS(&tls.Config{InsecureSkipVerify: true}),
		meilisearch.WithCustomRetries([]int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusInternalServerError, http.StatusServiceUnavailable}, uint8(maxRetries)),
	)
}
