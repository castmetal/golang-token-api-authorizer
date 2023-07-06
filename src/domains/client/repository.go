package client

import (
	"context"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/storage/postgres"
	"github.com/google/uuid"
)

type IClientRepository interface {
	common.IAggregateRoot
	Create(ctx context.Context, client *Client, qtx *postgres.Queries) error
	FindOneById(ctx context.Context, id uuid.UUID, qtx *postgres.Queries) (*Client, error)
	FindOneByName(ctx context.Context, name string, qtx *postgres.Queries) (*Client, error)
	GetActiveClientBySalt(ctx context.Context, salt string, id uuid.UUID, qtx *postgres.Queries) (*Client, error)
	GetActiveClientByApiId(ctx context.Context, apiId string, id uuid.UUID, qtx *postgres.Queries) (*Client, error)
	UpdateById(ctx context.Context, id uuid.UUID, client *Client, qtx *postgres.Queries) error
	DelById(ctx context.Context, id uuid.UUID, qtx *postgres.Queries) error
}
