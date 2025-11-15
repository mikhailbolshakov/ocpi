-- +goose Up
drop table parties;
create table parties
(
    id           varchar primary key,
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

-- +goose Down
drop table parties;
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