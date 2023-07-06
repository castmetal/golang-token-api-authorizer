package repositories

import (
	"context"
	"strings"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/common/logger"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/resource"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/storage/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ResourceRepository struct {
	db        *pgxpool.Pool
	pgQueries *postgres.Queries
}

func NewResourceRepository(db *pgxpool.Pool) resource.IResourceRepository {
	q := postgres.New(db)

	return &ResourceRepository{db: db, pgQueries: q}
}

func (r *ResourceRepository) Create(ctx context.Context, resource *resource.Resource, qtx *postgres.Queries) error {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	return q.CreateResource(ctx, postgres.CreateResourceParams{
		ID:                resource.ID,
		ResourceName:      resource.ResourceName,
		ResourcePath:      resource.ResourcePath,
		ResourceMethod:    postgres.RequestMethods(resource.ResourceMethod.String()),
		ResourceCreatedAt: resource.ResourceCreatedAt.Value,
		ResourceUpdatedAt: resource.ResourceUpdatedAt.Value,
	})
}

func (r *ResourceRepository) FindOneById(ctx context.Context, id uuid.UUID, qtx *postgres.Queries) (*resource.Resource, error) {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	rs, err := q.GetResourceById(ctx, id)
	if err != nil || rs.ID.String() == "" {
		return nil, common.NotFoundError("The resource id " + id.String())
	}

	mapper, err := r.ResourceMapper(rs)
	if err != nil {
		logger.Error(ctx, err, "client connection error")
		return nil, common.DefaultDomainError("client connection error")
	}

	return mapper, nil
}

func (r *ResourceRepository) FindOneByName(ctx context.Context, name string, qtx *postgres.Queries) (*resource.Resource, error) {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	rs, err := q.GetResourceByName(ctx, name)
	if err != nil || rs.ID.String() == "" {
		return nil, common.NotFoundError("The resource name " + name)
	}

	mapper, err := r.ResourceMapper(rs)
	if err != nil {
		logger.Error(ctx, err, "client connection error")
		return nil, common.DefaultDomainError("client connection error")
	}

	return mapper, nil
}

func (r *ResourceRepository) FindOneByPath(ctx context.Context, path string, method string, qtx *postgres.Queries) (*resource.Resource, error) {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	rs, err := q.GetResourceByPath(ctx, postgres.GetResourceByPathParams{
		ResourcePath:   path,
		ResourceMethod: postgres.RequestMethods(strings.ToUpper(method)),
	})
	if err != nil || rs.ID.String() == "" {
		return nil, common.NotFoundError("The resource path " + path)
	}

	mapper, err := r.ResourceMapper(rs)
	if err != nil {
		logger.Error(ctx, err, "client connection error")
		return nil, common.DefaultDomainError("client connection error")
	}

	return mapper, nil
}

func (r *ResourceRepository) DelById(ctx context.Context, id uuid.UUID, qtx *postgres.Queries) error {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	err := q.DelResourceById(ctx, id)
	if err != nil {
		logger.Error(ctx, err, "database off")
		return common.ConnectionClosedError("db_off")
	}

	return nil
}

func (r *ResourceRepository) ResourceMapper(data postgres.Resource) (*resource.Resource, error) {
	var rs *resource.Resource

	rs = &resource.Resource{
		ID:                data.ID,
		ResourceName:      data.ResourceName,
		ResourcePath:      data.ResourcePath,
		ResourceMethod:    resource.GetMethodByString(string(data.ResourceMethod)),
		ResourceCreatedAt: common.JsonTime{Value: data.ResourceCreatedAt},
		ResourceUpdatedAt: common.JsonTime{Value: data.ResourceUpdatedAt},
	}

	return rs, nil
}

func (r *ResourceRepository) GetDB() *pgxpool.Pool {
	return r.db
}

func (r *ResourceRepository) GetQueries() *postgres.Queries {
	return r.pgQueries
}
