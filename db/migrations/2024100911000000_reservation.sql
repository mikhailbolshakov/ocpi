-- +goose Up

alter table commands add reservation_id varchar GENERATED ALWAYS as (lower(details -> 'reserve' ->>  'reservationId')) stored;
create index idx_cmd_res on commands(reservation_id);

-- +goose Down
alter table commands drop reservation_id;

