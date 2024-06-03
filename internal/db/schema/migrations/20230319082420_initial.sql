-- +goose Up
-- +goose StatementBegin
CREATE TABLE translations
(
    id uuid PRIMARY KEY NOT NULL,
    key text UNIQUE NOT NULL,
    translation text NOT NULL,
    language_code text NOT NULL,
    CONSTRAINT translations_key_translation_language_code UNIQUE (key, translation, language_code)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS translations;
-- +goose StatementEnd
