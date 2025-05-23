package inf

import "github.com/meilisearch/meilisearch-go"

type IMeiliSearch interface {
	CreateCollection(name string, primaryKey string, schema any) error
	FindOne(name string, id string, filter *meilisearch.DocumentQuery, dest any) error
	Find(name string, filter *meilisearch.DocumentsQuery, dest *meilisearch.DocumentsResult) error
	Like(name string, query string, filter *meilisearch.SearchRequest, dest *meilisearch.SearchResponse) error
	Update(doc string, id string, value any) (*meilisearch.TaskInfo, error)
	BulkUpdate(doc string, id string, value any) (*meilisearch.TaskInfo, error)
	Delete(doc string, id string) (*meilisearch.TaskInfo, error)
	BulkDelete(doc string, ids ...string) (*meilisearch.TaskInfo, error)
}
