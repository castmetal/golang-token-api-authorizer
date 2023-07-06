-- name: CreateResource :exec
insert into resource (
    "id",
    "resource_name",
    "resource_path",
    "resource_method",
    "resource_created_at",
    "resource_updated_at"
) values ($1, $2, $3, $4, $5, $6);

-- name: GetResourceByName :one
select * from resource
where 
    resource_name = $1;

-- name: GetResourceById :one
select * from resource
where 
    id = $1;

-- name: DelResourceById :exec
delete from resource
where 
    id = $1;

-- name: GetResourceByPath :one
select * from resource
where 
    resource_path = $1 and resource_method = $2;

