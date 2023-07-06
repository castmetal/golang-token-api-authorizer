-- name: CreateScope :exec
insert into scope (
    "id",
    "scope_name",
    "scope_created_at",
    "scope_updated_at"
) values ($1, $2, $3, $4);

-- name: GetScopeByName :one
select * from scope
where 
    scope_name = $1;

-- name: GetScopeById :one
select * from scope
where 
    id = $1;

-- name: DelScopeById :exec
delete from scope
where 
    id = $1;
