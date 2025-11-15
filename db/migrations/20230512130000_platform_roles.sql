-- +goose Up

create table platforms
(
    id         varchar primary key,
    status     varchar   not null,
    role       varchar   not null,
    name       varchar,
    token_a    varchar   not null,
    token_b    varchar,
    token_c    varchar,
    details    jsonb,
    remote     bool      not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp
);

create index idx_platforms_token_a on platforms (token_a);
create index idx_platforms_token_b on platforms (token_b) where token_b is not null;
create index idx_platforms_token_c on platforms (token_c) where token_c is not null;

create table parties
(
    id           uuid primary key,
    platform_id  varchar   not null,
    party_id     varchar   not null,
    country_code varchar   not null,
    ref_id       varchar,
    roles        varchar[] not null,
    status       varchar   not null,
    details      jsonb,
    last_updated timestamp not null,
    last_sent    timestamp
);

create index idx_parties_roles on parties using gin ("roles");
create index idx_parties_platform on parties (platform_id);
create index idx_parties_party on parties (party_id, country_code);
create index idx_parties_ref on parties (ref_id);
create index idx_parties_last_upd on parties (last_updated);

create table logs
(
    event            varchar   not null,
    url              varchar,
    token            varchar,
    rq_id            varchar,
    corr_id          varchar,
    from_platform_id varchar,
    to_platform_id   varchar,
    rq_body          jsonb,
    rs_body          jsonb,
    rs_status        integer,
    ocpi_status      integer,
    err              text,
    created_at       timestamp not null
);

-- +goose Down
drop table platforms;
drop table parties;
drop table logs;
