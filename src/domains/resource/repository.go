package resource

import (
	"context"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/storage/postgres"
	"github.com/google/uuid"
)

type IResourceRepository interface {
	common.IAggregateRoot
	Create(ctx context.Context, resource *Resource, qtx *postgres.Queries) error
	FindOneById(ctx context.Context, id uuid.UUID, qtx *postgres.Queries) (*Resource, error)
	FindOneByName(ctx context.Context, name string, qtx *postgres.Queries) (*Resource, error)
	DelById(ctx context.Context, id uuid.UUID, qtx *postgres.Queries) error
	FindOneByPath(ctx context.Context, path string, method string, qtx *postgres.Queries) (*Resource, error)
}
