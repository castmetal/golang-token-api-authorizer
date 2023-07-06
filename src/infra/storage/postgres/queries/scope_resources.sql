-- name: CreateScopeResources :exec
insert into scope_resources (
    "id",
    "scope_id",
    "resource_id"
) values ($1, $2, $3);

-- name: GetScopeResourcesByIds :one
select * from scope_resources
where 
    scope_id = $1 and resource_id = $2;

-- name: GetScopeResourcesById :one
select * from scope_resources
where 
    id = $1;

-- name: DelScopeResourcesByIds :exec
delete from scope_resources
where 
    scope_id = $1 and resource_id = $2;
