CREATE TYPE period AS ENUM ('hours', 'minutes', 'seconds', 'days', 'years');

CREATE TABLE "client" (
  "id" uuid not null primary key,
  "client_name" varchar(60) NOT NULL,
  "scope_id" uuid NOT NULL,
  "permissions" jsonb NOT NULL,
  "api_id"  varchar(255) UNIQUE NOT NULL,
  "salt"  varchar(400) NOT NULL,
  "key_time_duration" decimal(3) NOT NULL,
  "key_period" period NOT NULL,
  "client_created_at" timestamptz not null default now(),
  "client_updated_at" timestamptz not null default now(),
  "client_deleted_at" timestamp
);

create index if not exists "client_scope_id_idx" on "client" (
    "scope_id"
);

create index if not exists "client_name_idx" on "client" (
    "client_name"
);

create index if not exists "client_api_id_idx" on "client" (
    "api_id"
);

create index if not exists "client_salt_idx" on "client" (
    "salt"
);

create index if not exists "client_created_at_idx" on "client" (
    "client_created_at"
);

create index if not exists "client_updated_at_idx" on "client" (
    "client_updated_at"
);

CREATE TABLE "scope" (
  "id" uuid not null primary key,
  "scope_name"  varchar(60) NOT NULL,
  "scope_created_at" timestamptz not null default now(),
  "scope_updated_at" timestamptz not null default now()
);

create index if not exists "scope_name_idx" on "scope" (
    "scope_name"
);

create index if not exists "scope_name_idx" on "scope" (
    "scope_created_at"
);

create index if not exists "scope_updated_at_idx" on "scope" (
    "scope_updated_at"
);

CREATE TABLE "scope_resources" (
  "id" uuid not null primary key,
  "scope_id" uuid NOT NULL,
  "resource_id" uuid NOT NULL
);

create index if not exists "scope_resources_scope_id_idx" on "scope_resources" (
    "scope_id"
);

create index if not exists "scope_resources_resource_id_idx" on "scope_resources" (
    "resource_id"
);

CREATE TYPE request_methods AS ENUM ('GET', 'POST', 'PUT', 'DELETE');

CREATE TABLE "resource" (
  "id" uuid not null primary key,
  "resource_name"  varchar(60) NOT NULL,
  "resource_path"  varchar(255) NOT NULL,
  "resource_method" request_methods NOT NULL,
  "resource_created_at" timestamptz not null default now(),
  "resource_updated_at" timestamptz not null default now()
);

create index if not exists "resource_name_idx" on "resource" (
    "resource_name"
);

create index if not exists "resource_created_at_idx" on "resource" (
    "resource_created_at"
);

create index if not exists "resource_updated_at_idx" on "resource" (
    "resource_updated_at"
);

ALTER TABLE "client" ADD FOREIGN KEY ("scope_id") REFERENCES "scope" ("id");

ALTER TABLE "scope_resources" ADD FOREIGN KEY ("scope_id") REFERENCES "scope" ("id");

ALTER TABLE "scope_resources" ADD FOREIGN KEY ("resource_id") REFERENCES "resource" ("id");
