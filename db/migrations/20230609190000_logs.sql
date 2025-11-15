-- +goose Up
alter table logs add incoming bool;
alter table logs add dur_ms integer;

-- +goose Down

alter table logs drop incoming;
alter table logs drop dur_ms;