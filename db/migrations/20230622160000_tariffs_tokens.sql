-- +goose Up

create table tariffs
(
    id           varchar primary key,
    platform_id  varchar   not null,
    party_id     varchar   not null,
    country_code varchar   not null,
    ref_id       varchar,
    details      jsonb,
    last_updated timestamp not null,
    last_sent    timestamp
);

create index idx_trf_platform on tariffs (platform_id);
create index idx_trf_party on tariffs (party_id, country_code);
create index idx_trf_ref on tariffs (ref_id);
create index idx_trf_last_upd on tariffs (last_updated);

create table tokens
(
    id           varchar primary key,
    platform_id  varchar   not null,
    party_id     varchar   not null,
    country_code varchar   not null,
    ref_id       varchar,
    details      jsonb,
    last_updated timestamp not null,
    last_sent    timestamp
);

create index idx_tkn_platform on tokens (platform_id);
create index idx_tkn_party on tokens (party_id, country_code);
create index idx_tkn_ref on tokens (ref_id);
create index idx_tkn_last_upd on tokens (last_updated);

-- +goose Down
drop table tariffs;
drop table tokens;
