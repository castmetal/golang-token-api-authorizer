// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: client.sql

package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

const createClient = `-- name: CreateClient :exec
insert into client (
    "id",
    "client_name",
    "scope_id",
    "permissions",
    "api_id",
    "salt",
    "key_time_duration",
    "key_period",
    "client_created_at",
    "client_updated_at"
) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

type CreateClientParams struct {
	ID              uuid.UUID
	ClientName      string
	ScopeID         uuid.UUID
	Permissions     pgtype.JSONB
	ApiID           string
	Salt            string
	KeyTimeDuration pgtype.Numeric
	KeyPeriod       Period
	ClientCreatedAt time.Time
	ClientUpdatedAt time.Time
}

func (q *Queries) CreateClient(ctx context.Context, arg CreateClientParams) error {
	_, err := q.db.Exec(ctx, createClient,
		arg.ID,
		arg.ClientName,
		arg.ScopeID,
		arg.Permissions,
		arg.ApiID,
		arg.Salt,
		arg.KeyTimeDuration,
		arg.KeyPeriod,
		arg.ClientCreatedAt,
		arg.ClientUpdatedAt,
	)
	return err
}

const delClientById = `-- name: DelClientById :exec
delete from client
where 
    client_deleted_at = $1
`

func (q *Queries) DelClientById(ctx context.Context, clientDeletedAt sql.NullTime) error {
	_, err := q.db.Exec(ctx, delClientById, clientDeletedAt)
	return err
}

const getActiveClientByApiId = `-- name: GetActiveClientByApiId :one
select id, client_name, scope_id, permissions, api_id, salt, key_time_duration, key_period, client_created_at, client_updated_at, client_deleted_at from client
where 
    api_id = $1 and id = $2 and client_deleted_at IS NULL
`

type GetActiveClientByApiIdParams struct {
	ApiID string
	ID    uuid.UUID
}

func (q *Queries) GetActiveClientByApiId(ctx context.Context, arg GetActiveClientByApiIdParams) (Client, error) {
	row := q.db.QueryRow(ctx, getActiveClientByApiId, arg.ApiID, arg.ID)
	var i Client
	err := row.Scan(
		&i.ID,
		&i.ClientName,
		&i.ScopeID,
		&i.Permissions,
		&i.ApiID,
		&i.Salt,
		&i.KeyTimeDuration,
		&i.KeyPeriod,
		&i.ClientCreatedAt,
		&i.ClientUpdatedAt,
		&i.ClientDeletedAt,
	)
	return i, err
}

const getActiveClientBySaltId = `-- name: GetActiveClientBySaltId :one
select id, client_name, scope_id, permissions, api_id, salt, key_time_duration, key_period, client_created_at, client_updated_at, client_deleted_at from client
where 
    salt = $1 and id = $2 and client_deleted_at IS NULL
`

type GetActiveClientBySaltIdParams struct {
	Salt string
	ID   uuid.UUID
}

func (q *Queries) GetActiveClientBySaltId(ctx context.Context, arg GetActiveClientBySaltIdParams) (Client, error) {
	row := q.db.QueryRow(ctx, getActiveClientBySaltId, arg.Salt, arg.ID)
	var i Client
	err := row.Scan(
		&i.ID,
		&i.ClientName,
		&i.ScopeID,
		&i.Permissions,
		&i.ApiID,
		&i.Salt,
		&i.KeyTimeDuration,
		&i.KeyPeriod,
		&i.ClientCreatedAt,
		&i.ClientUpdatedAt,
		&i.ClientDeletedAt,
	)
	return i, err
}

const getClientById = `-- name: GetClientById :one
select id, client_name, scope_id, permissions, api_id, salt, key_time_duration, key_period, client_created_at, client_updated_at, client_deleted_at from client
where 
    id = $1 and client_deleted_at IS NULL
`

func (q *Queries) GetClientById(ctx context.Context, id uuid.UUID) (Client, error) {
	row := q.db.QueryRow(ctx, getClientById, id)
	var i Client
	err := row.Scan(
		&i.ID,
		&i.ClientName,
		&i.ScopeID,
		&i.Permissions,
		&i.ApiID,
		&i.Salt,
		&i.KeyTimeDuration,
		&i.KeyPeriod,
		&i.ClientCreatedAt,
		&i.ClientUpdatedAt,
		&i.ClientDeletedAt,
	)
	return i, err
}

const getClientByName = `-- name: GetClientByName :one
select id, client_name, scope_id, permissions, api_id, salt, key_time_duration, key_period, client_created_at, client_updated_at, client_deleted_at from client
where 
    client_name = $1
`

func (q *Queries) GetClientByName(ctx context.Context, clientName string) (Client, error) {
	row := q.db.QueryRow(ctx, getClientByName, clientName)
	var i Client
	err := row.Scan(
		&i.ID,
		&i.ClientName,
		&i.ScopeID,
		&i.Permissions,
		&i.ApiID,
		&i.Salt,
		&i.KeyTimeDuration,
		&i.KeyPeriod,
		&i.ClientCreatedAt,
		&i.ClientUpdatedAt,
		&i.ClientDeletedAt,
	)
	return i, err
}

const updateClient = `-- name: UpdateClient :exec
update client 
SET
    "client_name" = $1,
    "scope_id" = $2,
    "permissions" = $3,
    "api_id" = $4,
    "salt" = $5,
    "key_time_duration" = $6,
    "key_period" = $7,
    "client_updated_at" = now()
WHERE 
    "id" = $8
`

type UpdateClientParams struct {
	ClientName      string
	ScopeID         uuid.UUID
	Permissions     pgtype.JSONB
	ApiID           string
	Salt            string
	KeyTimeDuration pgtype.Numeric
	KeyPeriod       Period
	ID              uuid.UUID
}

func (q *Queries) UpdateClient(ctx context.Context, arg UpdateClientParams) error {
	_, err := q.db.Exec(ctx, updateClient,
		arg.ClientName,
		arg.ScopeID,
		arg.Permissions,
		arg.ApiID,
		arg.Salt,
		arg.KeyTimeDuration,
		arg.KeyPeriod,
		arg.ID,
	)
	return err
}
