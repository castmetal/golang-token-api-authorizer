// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: resource.sql

package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createResource = `-- name: CreateResource :exec
insert into resource (
    "id",
    "resource_name",
    "resource_path",
    "resource_method",
    "resource_created_at",
    "resource_updated_at"
) values ($1, $2, $3, $4, $5, $6)
`

type CreateResourceParams struct {
	ID                uuid.UUID
	ResourceName      string
	ResourcePath      string
	ResourceMethod    RequestMethods
	ResourceCreatedAt time.Time
	ResourceUpdatedAt time.Time
}

func (q *Queries) CreateResource(ctx context.Context, arg CreateResourceParams) error {
	_, err := q.db.Exec(ctx, createResource,
		arg.ID,
		arg.ResourceName,
		arg.ResourcePath,
		arg.ResourceMethod,
		arg.ResourceCreatedAt,
		arg.ResourceUpdatedAt,
	)
	return err
}

const delResourceById = `-- name: DelResourceById :exec
delete from resource
where 
    id = $1
`

func (q *Queries) DelResourceById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, delResourceById, id)
	return err
}

const getResourceById = `-- name: GetResourceById :one
select id, resource_name, resource_path, resource_method, resource_created_at, resource_updated_at from resource
where 
    id = $1
`

func (q *Queries) GetResourceById(ctx context.Context, id uuid.UUID) (Resource, error) {
	row := q.db.QueryRow(ctx, getResourceById, id)
	var i Resource
	err := row.Scan(
		&i.ID,
		&i.ResourceName,
		&i.ResourcePath,
		&i.ResourceMethod,
		&i.ResourceCreatedAt,
		&i.ResourceUpdatedAt,
	)
	return i, err
}

const getResourceByName = `-- name: GetResourceByName :one
select id, resource_name, resource_path, resource_method, resource_created_at, resource_updated_at from resource
where 
    resource_name = $1
`

func (q *Queries) GetResourceByName(ctx context.Context, resourceName string) (Resource, error) {
	row := q.db.QueryRow(ctx, getResourceByName, resourceName)
	var i Resource
	err := row.Scan(
		&i.ID,
		&i.ResourceName,
		&i.ResourcePath,
		&i.ResourceMethod,
		&i.ResourceCreatedAt,
		&i.ResourceUpdatedAt,
	)
	return i, err
}

const getResourceByPath = `-- name: GetResourceByPath :one
select id, resource_name, resource_path, resource_method, resource_created_at, resource_updated_at from resource
where 
    resource_path = $1 and resource_method = $2
`

type GetResourceByPathParams struct {
	ResourcePath   string
	ResourceMethod RequestMethods
}

func (q *Queries) GetResourceByPath(ctx context.Context, arg GetResourceByPathParams) (Resource, error) {
	row := q.db.QueryRow(ctx, getResourceByPath, arg.ResourcePath, arg.ResourceMethod)
	var i Resource
	err := row.Scan(
		&i.ID,
		&i.ResourceName,
		&i.ResourcePath,
		&i.ResourceMethod,
		&i.ResourceCreatedAt,
		&i.ResourceUpdatedAt,
	)
	return i, err
}
