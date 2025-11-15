-- +goose Up

create table session_charging_periods
(
    session_id   varchar not null,
    details      jsonb,
    last_updated timestamp not null
);

create index idx_scp_session on session_charging_periods (session_id);

-- +goose Down
drop table session_charging_periods;
