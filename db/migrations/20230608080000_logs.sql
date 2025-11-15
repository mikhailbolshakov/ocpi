-- +goose Up
drop table logs;
create table logs
(
    event            varchar   not null,
    url              varchar,
    token            varchar,
    rq_id            varchar,
    corr_id          varchar,
    from_platform_id varchar,
    to_platform_id   varchar,
    details          jsonb,
    status           integer,
    ocpi_status      integer,
    err              text,
    created_at       timestamp not null
);

-- +goose Down
drop table logs;
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