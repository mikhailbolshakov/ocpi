-- +goose Up
alter table sessions add auth_ref varchar GENERATED ALWAYS as (details ->> 'authRef') stored;
create index idx_sess_auth_ref on sessions(auth_ref);

-- +goose Down
alter table sessions drop auth_ref;