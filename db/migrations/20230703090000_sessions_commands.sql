-- +goose Up

create table sessions
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

create index idx_sess_platform on sessions (platform_id);
create index idx_sess_party on sessions (party_id, country_code);
create index idx_sess_ref on sessions (ref_id);
create index idx_sess_last_upd on sessions (last_updated);

create table commands
(
    id           varchar primary key,
    status       varchar   not null,
    cmd          varchar   not null,
    deadline     timestamp,
    auth_ref     varchar,
    platform_id  varchar   not null,
    party_id     varchar   not null,
    country_code varchar   not null,
    details      jsonb,
    ref_id       varchar,
    last_updated timestamp not null,
    last_sent    timestamp
);

create index idx_cmd_platform on commands (platform_id);
create index idx_cmd_party on commands (party_id, country_code);
create index idx_cmd_ref on commands (ref_id);
create index idx_cmd_last_upd on commands (last_updated);
create index idx_cmd_deadline on commands (deadline);
create index idx_cmd_auth_ref on commands (auth_ref);

-- +goose Down
drop table sessions;
drop table commands;
