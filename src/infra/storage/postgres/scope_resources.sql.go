// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: scope_resources.sql

package postgres

import (
	"context"

	"github.com/google/uuid"
)

const createScopeResources = `-- name: CreateScopeResources :exec
insert into scope_resources (
    "id",
    "scope_id",
    "resource_id"
) values ($1, $2, $3)
`

type CreateScopeResourcesParams struct {
	ID         uuid.UUID
	ScopeID    uuid.UUID
	ResourceID uuid.UUID
}

func (q *Queries) CreateScopeResources(ctx context.Context, arg CreateScopeResourcesParams) error {
	_, err := q.db.Exec(ctx, createScopeResources, arg.ID, arg.ScopeID, arg.ResourceID)
	return err
}

const delScopeResourcesByIds = `-- name: DelScopeResourcesByIds :exec
delete from scope_resources
where 
    scope_id = $1 and resource_id = $2
`

type DelScopeResourcesByIdsParams struct {
	ScopeID    uuid.UUID
	ResourceID uuid.UUID
}

func (q *Queries) DelScopeResourcesByIds(ctx context.Context, arg DelScopeResourcesByIdsParams) error {
	_, err := q.db.Exec(ctx, delScopeResourcesByIds, arg.ScopeID, arg.ResourceID)
	return err
}

const getScopeResourcesById = `-- name: GetScopeResourcesById :one
select id, scope_id, resource_id from scope_resources
where 
    id = $1
`

func (q *Queries) GetScopeResourcesById(ctx context.Context, id uuid.UUID) (ScopeResource, error) {
	row := q.db.QueryRow(ctx, getScopeResourcesById, id)
	var i ScopeResource
	err := row.Scan(&i.ID, &i.ScopeID, &i.ResourceID)
	return i, err
}

const getScopeResourcesByIds = `-- name: GetScopeResourcesByIds :one
select id, scope_id, resource_id from scope_resources
where 
    scope_id = $1 and resource_id = $2
`

type GetScopeResourcesByIdsParams struct {
	ScopeID    uuid.UUID
	ResourceID uuid.UUID
}

func (q *Queries) GetScopeResourcesByIds(ctx context.Context, arg GetScopeResourcesByIdsParams) (ScopeResource, error) {
	row := q.db.QueryRow(ctx, getScopeResourcesByIds, arg.ScopeID, arg.ResourceID)
	var i ScopeResource
	err := row.Scan(&i.ID, &i.ScopeID, &i.ResourceID)
	return i, err
}
