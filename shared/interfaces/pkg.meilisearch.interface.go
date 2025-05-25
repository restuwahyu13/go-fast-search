package inf

import "github.com/meilisearch/meilisearch-go"

type IMeiliSearch interface {
	CreateCollection(name string, primaryKey string, schema any) error
	FindOne(doc string, id string, filter *meilisearch.DocumentQuery, dest any) error
	Find(doc string, filter *meilisearch.DocumentsQuery, dest *meilisearch.DocumentsResult) error
	Like(doc string, query string, filter *meilisearch.SearchRequest, dest *meilisearch.SearchResponse) error
	Insert(doc string, value any) (*meilisearch.TaskInfo, error)
	Update(doc string, id string, value any) (*meilisearch.TaskInfo, error)
	Delete(doc string, id string) (*meilisearch.TaskInfo, error)
	BulkInsert(doc string, value any) (*meilisearch.TaskInfo, error)
	BulkUpdate(doc string, value any) (*meilisearch.TaskInfo, error)
	BulkDelete(doc string, ids ...string) (*meilisearch.TaskInfo, error)
	UpdateFilterableAttributes(doc string, request []string) ([]string, error)
	UpdateSortableAttributes(doc string, request []string) ([]string, error)
	UpdateSearchableAttributes(doc string, request []string) ([]string, error)
	UpdateDisplayedAttributes(doc string, request []string) ([]string, error)
}
