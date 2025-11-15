-- +goose Up

alter table parties add created_at timestamp not null default now();
alter table parties add updated_at timestamp not null default now();
alter table parties add deleted_at timestamp;

alter table locations add created_at timestamp not null default now();
alter table locations add updated_at timestamp not null default now();
alter table locations add deleted_at timestamp;

alter table evses add created_at timestamp not null default now();
alter table evses add updated_at timestamp not null default now();
alter table evses add deleted_at timestamp;

alter table connectors add created_at timestamp not null default now();
alter table connectors add updated_at timestamp not null default now();
alter table connectors add deleted_at timestamp;

alter table tariffs add created_at timestamp not null default now();
alter table tariffs add updated_at timestamp not null default now();
alter table tariffs add deleted_at timestamp;

alter table tokens add created_at timestamp not null default now();
alter table tokens add updated_at timestamp not null default now();
alter table tokens add deleted_at timestamp;

alter table commands add created_at timestamp not null default now();
alter table commands add updated_at timestamp not null default now();
alter table commands add deleted_at timestamp;

alter table sessions add created_at timestamp not null default now();
alter table sessions add updated_at timestamp not null default now();
alter table sessions add deleted_at timestamp;

-- +goose Down

alter table parties drop created_at;
alter table parties drop updated_at;
alter table parties drop deleted_at;

alter table locations drop created_at;
alter table locations drop updated_at;
alter table locations drop deleted_at;

alter table evses drop created_at;
alter table evses drop updated_at;
alter table evses drop deleted_at;

alter table connectors drop created_at;
alter table connectors drop updated_at;
alter table connectors drop deleted_at;

alter table tariffs drop created_at;
alter table tariffs drop updated_at;
alter table tariffs drop deleted_at;

alter table tokens drop created_at;
alter table tokens drop updated_at;
alter table tokens drop deleted_at;

alter table commands drop created_at;
alter table commands drop updated_at;
alter table commands drop deleted_at;

alter table sessions drop created_at;
alter table sessions drop updated_at;
alter table sessions drop deleted_at;