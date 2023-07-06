package repositories

import (
	"context"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/common/logger"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/scope"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/storage/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ScopeRepository struct {
	db        *pgxpool.Pool
	pgQueries *postgres.Queries
}

func NewScopeRepository(db *pgxpool.Pool) scope.IScopeRepository {
	q := postgres.New(db)

	return &ScopeRepository{db: db, pgQueries: q}
}

func (r *ScopeRepository) Create(ctx context.Context, scope *scope.Scope, qtx *postgres.Queries) error {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	return q.CreateScope(ctx, postgres.CreateScopeParams{
		ID:             scope.ID,
		ScopeName:      scope.ScopeName,
		ScopeCreatedAt: scope.ScopeCreatedAt.Value,
		ScopeUpdatedAt: scope.ScopeUpdatedAt.Value,
	})
}

func (r *ScopeRepository) FindOneById(ctx context.Context, id uuid.UUID, qtx *postgres.Queries) (*scope.Scope, error) {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	sc, err := q.GetScopeById(ctx, id)
	if err != nil || sc.ID.String() == "" {
		return nil, common.NotFoundError("The resource id " + id.String())
	}

	mapper, err := r.ScopeMapper(sc)
	if err != nil {
		logger.Error(ctx, err, "client connection error")
		return nil, common.DefaultDomainError("client connection error")
	}

	return mapper, nil
}

func (r *ScopeRepository) FindOneByName(ctx context.Context, name string, qtx *postgres.Queries) (*scope.Scope, error) {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	sc, err := q.GetScopeByName(ctx, name)
	if err != nil || sc.ID.String() == "" {
		return nil, common.NotFoundError("The resource name " + name)
	}

	mapper, err := r.ScopeMapper(sc)
	if err != nil {
		logger.Error(ctx, err, "client connection error")
		return nil, common.DefaultDomainError("client connection error")
	}

	return mapper, nil
}

func (r *ScopeRepository) DelById(ctx context.Context, id uuid.UUID, qtx *postgres.Queries) error {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	err := q.DelScopeById(ctx, id)
	if err != nil {
		logger.Error(ctx, err, "database off")
		return common.ConnectionClosedError("db_off")
	}

	return nil
}

func (r *ScopeRepository) AttachResourceId(ctx context.Context, resourceId uuid.UUID, scopeId uuid.UUID, qtx *postgres.Queries) (uuid.UUID, error) {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	id := uuid.New()

	err := q.CreateScopeResources(ctx, postgres.CreateScopeResourcesParams{
		ID:         id,
		ScopeID:    scopeId,
		ResourceID: resourceId,
	})
	if err != nil {
		logger.Error(ctx, err, "database off")
		return uuid.Nil, common.ConnectionClosedError("db_off")
	}

	return id, nil
}

func (r *ScopeRepository) DetachResourceId(ctx context.Context, resourceId uuid.UUID, scopeId uuid.UUID, qtx *postgres.Queries) (uuid.UUID, error) {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	id := uuid.New()

	err := q.DelScopeResourcesByIds(ctx, postgres.DelScopeResourcesByIdsParams{
		ScopeID:    scopeId,
		ResourceID: resourceId,
	})
	if err != nil {
		logger.Error(ctx, err, "database off")
		return uuid.Nil, common.ConnectionClosedError("db_off")
	}

	return id, nil
}

func (r *ScopeRepository) ScopeMapper(data postgres.Scope) (*scope.Scope, error) {
	var sc *scope.Scope

	sc = &scope.Scope{
		ID:             data.ID,
		ScopeName:      data.ScopeName,
		ScopeCreatedAt: common.JsonTime{Value: data.ScopeCreatedAt},
		ScopeUpdatedAt: common.JsonTime{Value: data.ScopeUpdatedAt},
	}

	return sc, nil
}

func (r *ScopeRepository) GetDB() *pgxpool.Pool {
	return r.db
}

func (r *ScopeRepository) GetQueries() *postgres.Queries {
	return r.pgQueries
}
