package con

import (
	"crypto/tls"
	"net/http"

	"github.com/meilisearch/meilisearch-go"

	"github.com/restuwahyu13/go-fast-search/shared/dto"
)

func MeiliSearchConnection(req dto.Request[dto.Environtment]) meilisearch.ServiceManager {
	maxRetries := 15

	return meilisearch.New(req.Config.MEILISEARCH.URL,
		meilisearch.WithAPIKey(req.Config.MEILISEARCH.KEY),
		meilisearch.WithContentEncoding(meilisearch.BrotliEncoding, meilisearch.BestCompression),
		meilisearch.WithCustomRetries([]int{http.StatusInternalServerError, http.StatusServiceUnavailable}, uint8(maxRetries)),
		meilisearch.WithCustomClientWithTLS(&tls.Config{InsecureSkipVerify: true}),
	)
}
