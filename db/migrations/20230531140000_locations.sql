-- +goose Up

create table locations
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

create index idx_loc_platform on locations (platform_id);
create index idx_loc_party on locations (party_id, country_code);
create index idx_loc_ref on locations (ref_id);
create index idx_loc_last_upd on locations (last_updated);

create table evses
(
    id           varchar primary key,
    location_id  varchar,
    status       varchar,
    platform_id  varchar   not null,
    party_id     varchar   not null,
    country_code varchar   not null,
    ref_id       varchar,
    details      jsonb,
    last_updated timestamp not null,
    last_sent    timestamp
);

create index idx_evse_platform on evses (platform_id);
create index idx_evse_party on evses (party_id, country_code);
create index idx_evse_ref on evses (ref_id);
create index idx_evse_loc on evses (location_id);
create index idx_evse_last_upd on evses (last_updated);

create table connectors
(
    id           varchar primary key,
    location_id  varchar,
    evse_id      varchar,
    platform_id  varchar   not null,
    party_id     varchar   not null,
    country_code varchar   not null,
    ref_id       varchar,
    details      jsonb,
    last_updated timestamp not null,
    last_sent    timestamp
);

create index idx_con_platform on connectors (platform_id);
create index idx_con_party on connectors (party_id, country_code);
create index idx_con_ref on connectors (ref_id);
create index idx_con_loc on connectors (location_id);
create index idx_con_evse on connectors (evse_id);
create index idx_con_last_upd on connectors (last_updated);

-- +goose Down
drop table locations;
drop table evses;
drop table connectors;
