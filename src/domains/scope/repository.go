package scope

import (
	"context"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/storage/postgres"
	"github.com/google/uuid"
)

type IScopeRepository interface {
	common.IAggregateRoot
	Create(ctx context.Context, scope *Scope, qtx *postgres.Queries) error
	FindOneById(ctx context.Context, id uuid.UUID, qtx *postgres.Queries) (*Scope, error)
	FindOneByName(ctx context.Context, name string, qtx *postgres.Queries) (*Scope, error)
	DelById(ctx context.Context, id uuid.UUID, qtx *postgres.Queries) error
	AttachResourceId(ctx context.Context, resourceId uuid.UUID, scopeId uuid.UUID, qtx *postgres.Queries) (uuid.UUID, error)
	DetachResourceId(ctx context.Context, resourceId uuid.UUID, scopeId uuid.UUID, qtx *postgres.Queries) (uuid.UUID, error)
}
