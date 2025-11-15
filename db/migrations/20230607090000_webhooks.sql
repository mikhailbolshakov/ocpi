-- +goose Up

create table webhooks
(
    id         varchar primary key,
    api_key    varchar   not null,
    url        varchar   not null,
    events     varchar[] not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp
);

-- +goose Down
drop table webhooks;
