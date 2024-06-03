-- +goose Up
-- +goose StatementBegin
alter table translations
drop constraint translations_key_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table translations add constraint translations_key_key unique (key);
-- +goose StatementEnd
