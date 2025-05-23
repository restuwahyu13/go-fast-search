package pkg

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"time"

	search "github.com/meilisearch/meilisearch-go"

	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
)

type meilisearch struct {
	ctx         context.Context
	meilisearch search.ServiceManager
}

func NewMeiliSearch(ctx context.Context, con search.ServiceManager) inf.IMeiliSearch {
	return meilisearch{ctx: ctx, meilisearch: con}
}

func (p meilisearch) validate(doc string, value any) error {
	result, err := p.meilisearch.GetIndex(doc)
	if err != nil {
		return err
	}

	if result.UID != "" {
		return errors.New("doc: collection doc not exists in our system")
	}

	if value != nil && reflect.ValueOf(value).Kind() != reflect.Pointer {
		return errors.New("value must be a pointer")
	} else {
		if reflect.ValueOf(value).Elem().Kind() != reflect.Map || reflect.ValueOf(value).Elem().Kind() != reflect.Struct {
			return errors.New("value must be a map or struct")
		} else if reflect.ValueOf(value).Elem().Kind() != reflect.Slice {
			return errors.New("value must be a slice")
		}
	}

	return nil
}

func (p meilisearch) CreateCollection(name string, primaryKey string, schema any) error {
	result, err := p.meilisearch.GetIndexWithContext(p.ctx, name)
	if err != nil {
		return err
	}

	if result.UID == "" {
		if _, err := p.meilisearch.Index(name).AddDocumentsWithContext(p.ctx, schema, primaryKey); err != nil {
			return err
		}
	}

	return nil
}

func (p meilisearch) Insert(doc string, value any) (*search.TaskInfo, error) {
	if err := p.validate(doc, value); err != nil {
		return nil, err
	}

	res, err := p.meilisearch.Index(doc).AddDocumentsWithContext(p.ctx, value)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p meilisearch) FindOne(doc string, id string, filter *search.DocumentQuery, dest any) error {
	if err := p.validate(doc, nil); err != nil {
		return err
	}

	if err := p.meilisearch.Index(doc).GetDocumentWithContext(p.ctx, id, filter, dest); err != nil {
		return err
	}

	return nil
}

func (p meilisearch) Find(doc string, filter *search.DocumentsQuery, dest *search.DocumentsResult) error {
	if err := p.validate(doc, nil); err != nil {
		return err
	}

	if err := p.meilisearch.Index(doc).GetDocumentsWithContext(p.ctx, filter, dest); err != nil {
		return err
	}

	return nil
}

func (p meilisearch) Like(doc string, query string, filter *search.SearchRequest, dest *search.SearchResponse) error {
	if err := p.validate(doc, nil); err != nil {
		return err
	}

	dest, err := p.meilisearch.Index(doc).SearchWithContext(p.ctx, query, filter)
	if err != nil {
		return err
	}

	return nil
}

func (p meilisearch) Update(doc string, id string, value any) (*search.TaskInfo, error) {
	res := make(map[string]any)

	if err := p.validate(doc, value); err != nil {
		return nil, err
	}

	if err := p.FindOne(doc, id, &search.DocumentQuery{Fields: []string{"id"}}, &res); err != nil {
		return nil, err
	}

	if res == nil {
		return nil, sql.ErrNoRows
	}

	task, err := p.meilisearch.Index(doc).UpdateDocumentsWithContext(p.ctx, value)
	if err != nil {
		return nil, err
	}

	if task.TaskUID < 1 {
		return nil, cons.NO_ROWS_AFFECTED
	}

	if _, err := p.meilisearch.WaitForTaskWithContext(p.ctx, task.TaskUID, time.Duration(time.Second*3)); err != nil {
		return nil, err
	}

	return task, nil
}

func (p meilisearch) BulkUpdate(doc string, id string, value any) (*search.TaskInfo, error) {
	if err := p.validate(doc, value); err != nil {
		return nil, err
	}

	task, err := p.meilisearch.Index(doc).UpdateDocumentsWithContext(p.ctx, value)
	if err != nil {
		return nil, err
	}

	if task.TaskUID < 1 {
		return nil, cons.NO_ROWS_AFFECTED
	}

	if _, err := p.meilisearch.WaitForTaskWithContext(p.ctx, task.TaskUID, time.Duration(time.Second*3)); err != nil {
		return nil, err
	}

	return task, nil
}

func (p meilisearch) Delete(doc string, id string) (*search.TaskInfo, error) {
	res := make(map[string]any)

	if err := p.validate(doc, nil); err != nil {
		return nil, err
	}

	if err := p.FindOne(doc, id, &search.DocumentQuery{Fields: []string{"id"}}, &res); err != nil {
		return nil, err
	}

	if res == nil {
		return nil, sql.ErrNoRows
	}

	task, err := p.meilisearch.Index(doc).DeleteDocumentsWithContext(p.ctx, []string{id})
	if err != nil {
		return nil, err
	}

	if task.TaskUID < 1 {
		return nil, cons.NO_ROWS_AFFECTED
	}

	if _, err := p.meilisearch.WaitForTaskWithContext(p.ctx, task.TaskUID, time.Duration(time.Second*3)); err != nil {
		return nil, err
	}

	return task, nil
}

func (p meilisearch) BulkDelete(doc string, ids ...string) (*search.TaskInfo, error) {
	res := make(map[string]any)

	if err := p.validate(doc, nil); err != nil {
		return nil, err
	}

	for _, id := range ids {
		if err := p.FindOne(doc, id, &search.DocumentQuery{Fields: []string{"id"}}, &res); err != nil {
			return nil, err
		}

		if res == nil {
			return nil, sql.ErrNoRows
		}
	}

	task, err := p.meilisearch.Index(doc).DeleteDocumentsWithContext(p.ctx, ids)
	if err != nil {
		return nil, err
	}

	if task.TaskUID < 1 {
		return nil, cons.NO_ROWS_AFFECTED
	}

	if _, err := p.meilisearch.WaitForTaskWithContext(p.ctx, task.TaskUID, time.Duration(time.Second*3)); err != nil {
		return nil, err
	}

	return task, nil
}
