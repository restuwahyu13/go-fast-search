package pkg

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"slices"
	"time"

	search "github.com/meilisearch/meilisearch-go"

	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	helper "github.com/restuwahyu13/go-fast-search/shared/helpers"
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

	if result == nil {
		return errors.New("doc: collection not exists in our system")
	}

	valueof := reflect.ValueOf(value)
	if value != nil && valueof.Kind() == reflect.Pointer {
		elemof := valueof.Elem().Kind()
		validElem := []reflect.Kind{reflect.Map, reflect.Struct, reflect.Slice}

		if slices.Index(validElem, elemof) == -1 {
			return errors.New("value must be a map, struct or slice")
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

func (p meilisearch) BulkInsert(doc string, value any) (*search.TaskInfo, error) {
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

	result, err := p.meilisearch.Index(doc).SearchWithContext(p.ctx, query, filter)
	if err != nil {
		return err
	}

	transform := helper.NewTransform()
	if err := transform.SrcToDest(result, dest); err != nil {
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

func (p meilisearch) BulkUpdate(doc string, value any) (*search.TaskInfo, error) {
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

	if err := p.FindOne(doc, id, &search.DocumentQuery{}, &res); err != nil {
		return nil, err
	}

	if res == nil {
		return nil, sql.ErrNoRows
	}

	res["deleted_at"] = time.Now().Format(cons.DATE_TIME_FORMAT)

	task, err := p.meilisearch.Index(doc).UpdateDocumentsWithContext(p.ctx, &res)
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
	resDoc := make(map[string]any)
	resDcos := []map[string]any{}

	if err := p.validate(doc, nil); err != nil {
		return nil, err
	}

	for _, id := range ids {
		if err := p.FindOne(doc, id, &search.DocumentQuery{}, &resDoc); err != nil {
			return nil, err
		}

		if resDoc == nil {
			return nil, sql.ErrNoRows
		}

		resDoc["deleted_at"] = time.Now().Format(cons.DATE_TIME_FORMAT)
		resDcos = append(resDcos, resDoc)
	}

	task, err := p.meilisearch.Index(doc).UpdateDocumentsWithContext(p.ctx, &resDcos)
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

func (p meilisearch) UpdateFilterableAttributes(doc string, request []string) ([]string, error) {
	filterAttributesPtr, err := p.meilisearch.Index(doc).GetFilterableAttributesWithContext(p.ctx)
	if err != nil {
		return nil, err
	}

	filterAbleIdx := 0
	for _, attribute := range request {
		if slices.Index(*filterAttributesPtr, attribute) == -1 {
			filterAbleIdx = -1
		}
	}

	if filterAbleIdx == -1 {
		task, err := p.meilisearch.Index(doc).UpdateFilterableAttributesWithContext(p.ctx, &request)
		if err != nil {
			return nil, err
		}

		if task.TaskUID < 1 {
			return nil, cons.NO_ROWS_AFFECTED
		}

		if _, err := p.meilisearch.WaitForTaskWithContext(p.ctx, task.TaskUID, time.Duration(time.Second*3)); err != nil {
			return nil, err
		}

	}

	return *filterAttributesPtr, nil
}

func (p meilisearch) UpdateSortableAttributes(doc string, request []string) ([]string, error) {
	sortAttributesPtr, err := p.meilisearch.Index(doc).GetSortableAttributesWithContext(p.ctx)
	if err != nil {
		return nil, err
	}

	sortAbleIdx := 0
	for _, attribute := range request {
		if slices.Index(*sortAttributesPtr, attribute) == -1 {
			sortAbleIdx = -1
		}
	}

	if sortAbleIdx == -1 {
		task, err := p.meilisearch.Index(doc).UpdateSortableAttributesWithContext(p.ctx, &request)
		if err != nil {
			return nil, err
		}

		if task.TaskUID < 1 {
			return nil, cons.NO_ROWS_AFFECTED
		}

		if _, err := p.meilisearch.WaitForTaskWithContext(p.ctx, task.TaskUID, time.Duration(time.Second*3)); err != nil {
			return nil, err
		}

		return request, nil
	}

	return *sortAttributesPtr, nil
}

func (p meilisearch) UpdateSearchableAttributes(doc string, request []string) ([]string, error) {
	searchAttributesPtr, err := p.meilisearch.Index(doc).GetSearchableAttributesWithContext(p.ctx)
	if err != nil {
		return nil, err
	}

	sortAbleIdx := 0
	for _, attribute := range request {
		if slices.Index(*searchAttributesPtr, attribute) == -1 {
			sortAbleIdx = -1
		}
	}

	if sortAbleIdx == -1 {
		task, err := p.meilisearch.Index(doc).UpdateSearchableAttributesWithContext(p.ctx, &request)
		if err != nil {
			return nil, err
		}

		if task.TaskUID < 1 {
			return nil, cons.NO_ROWS_AFFECTED
		}

		if _, err := p.meilisearch.WaitForTaskWithContext(p.ctx, task.TaskUID, time.Duration(time.Second*3)); err != nil {
			return nil, err
		}

		return request, nil
	}

	return *searchAttributesPtr, nil
}

func (p meilisearch) UpdateDisplayedAttributes(doc string, request []string) ([]string, error) {
	searchAttributesPtr, err := p.meilisearch.Index(doc).GetDisplayedAttributesWithContext(p.ctx)
	if err != nil {
		return nil, err
	}

	sortAbleIdx := 0
	for _, attribute := range request {
		if slices.Index(*searchAttributesPtr, attribute) == -1 {
			sortAbleIdx = -1
		}
	}

	if sortAbleIdx == -1 {
		task, err := p.meilisearch.Index(doc).UpdateDisplayedAttributesWithContext(p.ctx, &request)
		if err != nil {
			return nil, err
		}

		if task.TaskUID < 1 {
			return nil, cons.NO_ROWS_AFFECTED
		}

		if _, err := p.meilisearch.WaitForTaskWithContext(p.ctx, task.TaskUID, time.Duration(time.Second*3)); err != nil {
			return nil, err
		}

		return request, nil
	}

	return *searchAttributesPtr, nil
}
