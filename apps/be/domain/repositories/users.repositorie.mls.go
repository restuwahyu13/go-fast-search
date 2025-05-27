package repo

import (
	"context"
	"fmt"
	"math"

	"github.com/meilisearch/meilisearch-go"

	entitie "github.com/restuwahyu13/go-fast-search/domain/entities"
	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	helper "github.com/restuwahyu13/go-fast-search/shared/helpers"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
	"github.com/restuwahyu13/go-fast-search/shared/pkg"
)

type usersMeilisearchRepositorie struct {
	ctx         context.Context
	meilisearch inf.IMeiliSearch
	doc         *entitie.UsersDocument
}

func NewUsersMeilisearchRepositorie(ctx context.Context, db meilisearch.ServiceManager) inf.IUsersMeiliSearchRepositorie {
	meilisearch := pkg.NewMeiliSearch(ctx, db)
	return usersMeilisearchRepositorie{ctx: ctx, meilisearch: meilisearch, doc: new(entitie.UsersDocument)}
}

func (r usersMeilisearchRepositorie) Search(query string, filter *meilisearch.SearchRequest) (*opt.MeiliSearchDocuments[[]entitie.UsersDocument], error) {
	transform := helper.NewTransform()

	docResult := new(meilisearch.SearchResponse)
	docResultReformat := new(opt.MeiliSearchDocuments[[]entitie.UsersDocument])

	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return nil, err
	}

	err := r.meilisearch.Like("users", query, filter, docResult)
	if err != nil {
		return nil, err
	}

	if err := transform.SrcToDest(docResult, docResultReformat); err != nil {
		return nil, err
	}

	return docResultReformat, nil
}

func (r usersMeilisearchRepositorie) Find(filter *meilisearch.DocumentsQuery) (*opt.MeiliSearchDocuments[[]entitie.UsersDocument], error) {
	transform := helper.NewTransform()

	docResult := new(meilisearch.DocumentsResult)
	docResultReformat := new(opt.MeiliSearchDocuments[[]entitie.UsersDocument])

	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return nil, err
	}

	err := r.meilisearch.Find("users", filter, docResult)
	if err != nil {
		return nil, err
	}

	if err := transform.SrcToDest(docResult, docResultReformat); err != nil {
		return nil, err
	}

	return docResultReformat, nil
}

func (r usersMeilisearchRepositorie) FindOne(id string, filter *meilisearch.DocumentQuery) (*entitie.UsersDocument, error) {
	res := new(entitie.UsersDocument)

	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return nil, err
	}

	err := r.meilisearch.FindOne("users", id, filter, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r usersMeilisearchRepositorie) Insert(value any) error {
	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return err
	}

	if _, err := r.meilisearch.Insert("users", value); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) Update(id string, value any) error {
	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return err
	}

	if _, err := r.meilisearch.Update("users", id, value); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) Delete(id string) error {
	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return err
	}

	if _, err := r.meilisearch.Delete("users", id); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) BulkInsert(value any) error {
	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return err
	}

	if _, err := r.meilisearch.BulkInsert("users", value); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) BulkUpdate(value any) error {
	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return err
	}

	if _, err := r.meilisearch.BulkUpdate("users", value); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) BulkDelete(ids ...string) error {
	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return err
	}

	if _, err := r.meilisearch.BulkDelete("users", ids...); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) UpdateFilterableAttributes(attributes ...string) error {
	if _, err := r.meilisearch.UpdateFilterableAttributes("users", attributes); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) UpdateSearchableAttributes(attributes ...string) error {
	if _, err := r.meilisearch.UpdateSearchableAttributes("users", attributes); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) UpdateSortableAttributes(attributes ...string) error {
	if _, err := r.meilisearch.UpdateSortableAttributes("users", attributes); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) UpdateDisplayedAttributes(attributes ...string) error {
	if _, err := r.meilisearch.UpdateDisplayedAttributes("users", attributes); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) ListUsersDocuments(req dto.Request[dto.MeiliSearchDocumentsQuery]) (*opt.MeiliSearchDocuments[[]entitie.UsersDocument], error) {
	transform := helper.NewTransform()

	usersDocumentsResult := new(opt.MeiliSearchDocuments[[]entitie.UsersDocument])
	fields := []string{
		"id",
		"name",
		"email",
		"phone",
		"date_of_birth",
		"age",
		"address",
		"city",
		"state",
		"direction",
		"country",
		"postal_code",
		"created_at",
	}

	usersFilterDoc := new(dto.ListUsersFilterDTO)
	if err := transform.ReqToRes(&req.Query.Filter, usersFilterDoc); err != nil {
		return nil, err
	}

	filter := "deleted_at IS NULL"

	if usersFilterDoc.StartDate != "" && usersFilterDoc.EndDate != "" {
		startDate, err := helper.TimeStampToUnix(usersFilterDoc.StartDate)
		if err != nil {
			return nil, err
		}

		endDate, err := helper.TimeStampToUnix(usersFilterDoc.EndDate)
		if err != nil {
			return nil, err
		}

		filter += fmt.Sprintf(" AND created_at > %d AND created_at < %d", startDate, endDate)
	}

	if usersFilterDoc.Age != "" {
		filter += fmt.Sprintf(" AND age = %s", usersFilterDoc.Age)
	}

	if usersFilterDoc.City != "" {
		filter += fmt.Sprintf(" AND city = %s", usersFilterDoc.City)
	}

	if usersFilterDoc.State != "" {
		filter += fmt.Sprintf(" AND state = %s", usersFilterDoc.State)
	}

	if usersFilterDoc.Direction != "" {
		filter += fmt.Sprintf(" AND direction = %s", usersFilterDoc.Direction)
	}

	if usersFilterDoc.Country != "" {
		filter += fmt.Sprintf(" AND country = %s", usersFilterDoc.Country)
	}

	/**
	* FETCH DATA TERITORY
	 */

	if req.Query.Search == "" {
		mlsFetchReq := new(meilisearch.DocumentsQuery)
		mlsFetchReq.Limit = req.Query.Limit
		mlsFetchReq.Offset = req.Query.Page
		mlsFetchReq.Filter = filter
		mlsFetchReq.Fields = fields

		mlsFetchDcouments, err := r.Find(mlsFetchReq)
		if err != nil {
			return nil, err
		}

		usersDocumentsResult.Results = mlsFetchDcouments.Results
		usersDocumentsResult.Query = req.Query.Search
		usersDocumentsResult.Limit = req.Query.Limit
		usersDocumentsResult.Offset = req.Query.Page + 1
		usersDocumentsResult.TotalPages = int64(math.Ceil(float64(mlsFetchDcouments.Total) / float64(usersDocumentsResult.Limit)))
		usersDocumentsResult.Total = mlsFetchDcouments.Total
	}

	/**
	* SEARCH DATA TERITORY
	 */

	if req.Query.Search != "" {
		mlsSearchReq := new(meilisearch.SearchRequest)
		mlsSearchReq.Limit = req.Query.Limit
		mlsSearchReq.HitsPerPage = req.Query.Limit
		mlsSearchReq.Page = req.Query.Page
		mlsSearchReq.ShowMatchesPosition = cons.TRUE
		mlsSearchReq.AttributesToRetrieve = fields
		mlsSearchReq.Filter = filter

		if req.Query.MatchingStrategy != "" {
			switch req.Query.MatchingStrategy {

			case string(meilisearch.Last):
				mlsSearchReq.MatchingStrategy = meilisearch.Last
				break

			case string(meilisearch.All):
				mlsSearchReq.MatchingStrategy = meilisearch.All
				break

			case string(meilisearch.Frequency):
				mlsSearchReq.MatchingStrategy = meilisearch.Frequency
				break
			}
		}

		usersStatDocuments, err := r.meilisearch.GetStats("users")
		if err != nil {
			return nil, err
		}

		mlsSearchReq.AttributesToHighlight, err = r.meilisearch.GetSearchableAttributes("users")
		if err != nil {
			return nil, err
		}

		usersSearchDocuments, err := r.Search(req.Query.Search, mlsSearchReq)
		if err != nil {
			return nil, err
		}

		usersDocumentsResult.Results = usersSearchDocuments.Hits
		usersDocumentsResult.Query = req.Query.Search
		usersDocumentsResult.Limit = req.Query.Limit
		usersDocumentsResult.Offset = req.Query.Page + 1
		usersDocumentsResult.TotalPages = int64(math.Ceil(float64(usersStatDocuments.NumberOfDocuments) / float64(usersDocumentsResult.Limit)))
		usersDocumentsResult.Total = usersStatDocuments.NumberOfDocuments
	}

	return usersDocumentsResult, nil
}
