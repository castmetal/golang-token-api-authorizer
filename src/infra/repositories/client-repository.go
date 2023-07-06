package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math/big"
	"time"

	"github.com/castmetal/golang-token-api-authorizer/src/domains/client"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/common"
	"github.com/castmetal/golang-token-api-authorizer/src/domains/common/logger"
	"github.com/castmetal/golang-token-api-authorizer/src/infra/storage/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ClientRepository struct {
	db        *pgxpool.Pool
	pgQueries *postgres.Queries
}

func NewClientRepository(db *pgxpool.Pool) client.IClientRepository {
	q := postgres.New(db)

	return &ClientRepository{db: db, pgQueries: q}
}

func (r *ClientRepository) Create(ctx context.Context, client *client.Client, qtx *postgres.Queries) error {
	keyDuration := big.NewInt(int64(client.KeyTimeDuration))
	var permissions pgtype.JSONB

	permissionsBytes, err := json.Marshal(client.Permissions)
	if err != nil {
		return err
	}

	if err := permissions.Set(permissionsBytes); err != nil {
		return err
	}

	q := GetQuerieTransaction(qtx, r.pgQueries)

	return q.CreateClient(ctx, postgres.CreateClientParams{
		ID:              client.ID,
		ClientName:      client.ClientName,
		ScopeID:         client.ScopeId,
		Permissions:     permissions,
		ApiID:           client.ApiId,
		Salt:            client.Salt,
		KeyTimeDuration: pgtype.Numeric{Int: keyDuration, Status: pgtype.Present},
		KeyPeriod:       postgres.Period(client.KeyPeriod),
		ClientCreatedAt: client.ClientCreatedAt.Value,
		ClientUpdatedAt: client.ClientUpdatedAt.Value,
	})
}

func (r *ClientRepository) FindOneById(ctx context.Context, id uuid.UUID, qtx *postgres.Queries) (*client.Client, error) {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	clnt, err := q.GetClientById(ctx, id)
	if err != nil || clnt.ID.String() == "" {
		return nil, common.NotFoundError("The client id " + id.String())
	}

	mapper, err := r.ClientMapper(clnt)
	if err != nil {
		logger.Error(ctx, err, "client connection error")
		return nil, common.DefaultDomainError("client connection error")
	}

	return mapper, nil
}

func (r *ClientRepository) FindOneByName(ctx context.Context, name string, qtx *postgres.Queries) (*client.Client, error) {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	clnt, err := q.GetClientByName(ctx, name)
	if err != nil || clnt.ID.String() == "" {
		return nil, common.NotFoundError("The client: name " + name)
	}

	mapper, err := r.ClientMapper(clnt)
	if err != nil {
		logger.Error(ctx, err, "client connection error")
		return nil, common.DefaultDomainError("client connection error")
	}

	return mapper, nil
}

func (r *ClientRepository) GetActiveClientBySalt(ctx context.Context, salt string, id uuid.UUID, qtx *postgres.Queries) (*client.Client, error) {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	clnt, err := q.GetActiveClientBySaltId(ctx, postgres.GetActiveClientBySaltIdParams{
		Salt: salt,
		ID:   id,
	})
	if err != nil || clnt.ID.String() == "" {
		return nil, common.NotFoundError("The client: id " + id.String())
	}

	mapper, err := r.ClientMapper(clnt)
	if err != nil {
		logger.Error(ctx, err, "client connection error")
		return nil, common.DefaultDomainError("client connection error")
	}

	return mapper, nil
}

func (r *ClientRepository) GetActiveClientByApiId(ctx context.Context, apiId string, id uuid.UUID, qtx *postgres.Queries) (*client.Client, error) {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	clnt, err := q.GetActiveClientByApiId(ctx, postgres.GetActiveClientByApiIdParams{
		ApiID: apiId,
		ID:    id,
	})
	if err != nil || clnt.ID.String() == "" {
		return nil, common.NotFoundError("The client: api_id " + apiId)
	}

	mapper, err := r.ClientMapper(clnt)
	if err != nil {
		logger.Error(ctx, err, "client connection error")
		return nil, common.DefaultDomainError("client connection error")
	}

	return mapper, nil
}

func (r *ClientRepository) UpdateById(ctx context.Context, id uuid.UUID, client *client.Client, qtx *postgres.Queries) error {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	err := q.UpdateClient(ctx, postgres.UpdateClientParams{
		ClientName: client.ClientName,
	})
	if err != nil {
		logger.Error(ctx, err, "database off")
		return common.ConnectionClosedError("db_off")
	}

	return nil
}

func (r *ClientRepository) DelById(ctx context.Context, id uuid.UUID, qtx *postgres.Queries) error {
	q := GetQuerieTransaction(qtx, r.pgQueries)

	actualDate := time.Now()
	err := q.DelClientById(ctx, sql.NullTime{Time: actualDate, Valid: true})
	if err != nil {
		logger.Error(ctx, err, "database off")
		return common.ConnectionClosedError("db_off")
	}

	return nil
}

func (r *ClientRepository) ClientMapper(data postgres.Client) (*client.Client, error) {
	var clnt *client.Client

	var permissions []client.Permission

	permissionsBytes, err := json.Marshal(data.Permissions)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(permissionsBytes, &permissions); err != nil {
		log.Fatal(err)
	}

	clnt = &client.Client{
		ID:              data.ID,
		ClientName:      data.ClientName,
		ScopeId:         data.ScopeID,
		Permissions:     permissions,
		ApiId:           data.ApiID,
		Salt:            data.Salt,
		KeyTimeDuration: int32(data.KeyTimeDuration.Int.Int64()),
		KeyPeriod:       string(data.KeyPeriod),
		ClientCreatedAt: common.JsonTime{Value: data.ClientCreatedAt},
		ClientUpdatedAt: common.JsonTime{Value: data.ClientUpdatedAt},
	}

	if data.ClientDeletedAt.Valid && data.ClientDeletedAt.Time.String() != "" {
		clnt.ClientDeletedAt = &common.JsonNullTime{Value: data.ClientDeletedAt}
	}

	return clnt, nil
}

func (r *ClientRepository) GetDB() *pgxpool.Pool {
	return r.db
}

func (r *ClientRepository) GetQueries() *postgres.Queries {
	return r.pgQueries
}
