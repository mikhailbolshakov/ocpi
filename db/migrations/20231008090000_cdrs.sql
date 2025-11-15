-- +goose Up

create table cdrs
(
    id           varchar primary key,
    session_id   varchar,
    platform_id  varchar   not null,
    party_id     varchar   not null,
    country_code varchar   not null,
    ref_id       varchar,
    details      jsonb,
    last_updated timestamp not null,
    last_sent    timestamp,
    created_at   timestamp not null,
    updated_at   timestamp not null,
    deleted_at   timestamp
);

create index idx_cdr_platform on cdrs (platform_id);
create index idx_cdr_party on cdrs (party_id, country_code);
create index idx_cdr_ref on cdrs (ref_id) where ref_id is not null;
create index idx_cdr_last_upd on cdrs (last_updated);
create index idx_cdr_sess on cdrs (session_id);

-- +goose Down
drop table cdrs;
