package models

import (
	"context"
	"errors"
	"mortenvistisen/internal/storage"
	"mortenvistisen/internal/validation"

	"github.com/uptrace/bun"
)

type ArticleTagConnectionEntity struct {
	bun.BaseModel `      bun:"table:article_tag_connections,alias:article_tag_connections"`
	ID            int32 `bun:"id,pk,autoincrement"`
	ArticleID     int32 `bun:"article_id,autoincrement"`
	TagID         int32 `bun:"tag_id,autoincrement"`
}

func (e *ArticleTagConnectionEntity) Validate() error {
	return nil
}

func (atc articleTagConnection) Find(
	ctx context.Context,
	db storage.Executor,
	id int32,
) (ArticleTagConnectionEntity, error) {
	var entity ArticleTagConnectionEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return ArticleTagConnectionEntity{}, err
	}

	return entity, nil
}

type CreateArticleTagConnectionData struct {
	ArticleID int32
	TagID     int32
}

func (atc articleTagConnection) Create(
	ctx context.Context,
	db storage.Executor,
	data CreateArticleTagConnectionData,
) (ArticleTagConnectionEntity, error) {
	entity := ArticleTagConnectionEntity{
		ArticleID: data.ArticleID,
		TagID:     data.TagID,
	}

	if err := validation.Validate(&entity); err != nil {
		return ArticleTagConnectionEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if _, err := db.NewInsert().Model(&entity).Exec(ctx); err != nil {
		return ArticleTagConnectionEntity{}, err
	}

	return entity, nil
}

type UpdateArticleTagConnectionData struct {
	ID        int32
	ArticleID int32
	TagID     int32
}

func (atc articleTagConnection) Update(
	ctx context.Context,
	db storage.Executor,
	data UpdateArticleTagConnectionData,
) (ArticleTagConnectionEntity, error) {
	entity := ArticleTagConnectionEntity{
		ID:        data.ID,
		ArticleID: data.ArticleID,
		TagID:     data.TagID,
	}

	if err := validation.Validate(&entity); err != nil {
		return ArticleTagConnectionEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewUpdate().
		Model(&entity).
		Column("article_id").
		Column("tag_id").
		WherePK().
		Returning("*").
		Scan(ctx); err != nil {
		return ArticleTagConnectionEntity{}, err
	}

	return entity, nil
}

func (atc articleTagConnection) Destroy(ctx context.Context, db storage.Executor, id int32) error {
	_, err := db.NewDelete().
		Model((*ArticleTagConnectionEntity)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (atc articleTagConnection) All(
	ctx context.Context,
	db storage.Executor,
) ([]ArticleTagConnectionEntity, error) {
	var entities []ArticleTagConnectionEntity
	if err := db.NewSelect().
		Model(&entities).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

func (atc articleTagConnection) TagsForArticle(
	ctx context.Context,
	db storage.Executor,
	articleID int32,
) ([]TagEntity, error) {
	tags := make([]TagEntity, 0)
	if err := db.NewSelect().
		Model(&tags).
		Join("JOIN article_tag_connections AS atc ON atc.tag_id = tags.id").
		Where("atc.article_id = ?", articleID).
		OrderExpr("LOWER(tags.title) ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return tags, nil
}

func (atc articleTagConnection) ReplaceForArticle(
	ctx context.Context,
	db storage.Executor,
	articleID int32,
	tagIDs []int32,
) error {
	uniqueTagIDs := uniquePositiveIDs(tagIDs)
	deleteQuery := db.NewDelete().
		Model((*ArticleTagConnectionEntity)(nil)).
		Where("article_id = ?", articleID)
	if len(uniqueTagIDs) > 0 {
		deleteQuery = deleteQuery.Where("tag_id NOT IN (?)", bun.In(uniqueTagIDs))
	}
	if _, err := deleteQuery.Exec(ctx); err != nil {
		return err
	}

	if len(uniqueTagIDs) == 0 {
		return nil
	}

	connections := make([]ArticleTagConnectionEntity, 0, len(uniqueTagIDs))
	for _, tagID := range uniqueTagIDs {
		connections = append(connections, ArticleTagConnectionEntity{
			ArticleID: articleID,
			TagID:     tagID,
		})
	}

	_, err := db.NewInsert().
		Model(&connections).
		On("CONFLICT (article_id, tag_id) DO NOTHING").
		Exec(ctx)

	return err
}

func uniquePositiveIDs(ids []int32) []int32 {
	unique := make([]int32, 0, len(ids))
	seen := make(map[int32]struct{}, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}

	return unique
}

type PaginatedArticleTagConnections struct {
	ArticleTagConnections []ArticleTagConnectionEntity
	TotalCount            int64
	Page                  int64
	PageSize              int64
	TotalPages            int64
}

func (atc articleTagConnection) Paginate(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
) (PaginatedArticleTagConnections, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	totalCount, err := db.NewSelect().
		Model(&ArticleTagConnectionEntity{}).Count(ctx)
	if err != nil {
		return PaginatedArticleTagConnections{}, err
	}

	entities := make([]ArticleTagConnectionEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedArticleTagConnections{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedArticleTagConnections{
		ArticleTagConnections: entities,
		TotalCount:            int64(totalCount),
		Page:                  page,
		PageSize:              pageSize,
		TotalPages:            totalPages,
	}, nil
}

func (atc articleTagConnection) Upsert(
	ctx context.Context,
	db storage.Executor,
	data CreateArticleTagConnectionData,
) (ArticleTagConnectionEntity, error) {
	entity := ArticleTagConnectionEntity{
		ArticleID: data.ArticleID,
		TagID:     data.TagID,
	}

	if err := validation.Validate(&entity); err != nil {
		return ArticleTagConnectionEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (id) DO UPDATE").
		Set("article_id = excluded.article_id").
		Set("tag_id = excluded.tag_id").
		Returning("*").
		Scan(ctx); err != nil {
		return ArticleTagConnectionEntity{}, err
	}

	return entity, nil
}
