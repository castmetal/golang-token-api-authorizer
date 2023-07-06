-- name: CreateClient :exec
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
) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetClientByName :one
select * from client
where 
    client_name = $1;

-- name: GetActiveClientByApiId :one
select * from client
where 
    api_id = $1 and id = $2 and client_deleted_at IS NOT NULL;

-- name: GetActiveClientBySaltId :one
select * from client
where 
    salt = $1 and id = $2 and client_deleted_at IS NOT NULL;

-- name: GetClientById :one
select * from client
where 
    id = $1 and client_deleted_at IS NOT NULL;

-- name: DelClientById :exec
delete from client
where 
    client_deleted_at = $1;
    
-- name: UpdateClient :exec
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
    "id" = $8;
