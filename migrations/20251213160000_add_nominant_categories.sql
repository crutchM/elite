-- +goose Up
-- +goose StatementBegin

-- Удаляем внешний ключ и колонку category_id из nominants
ALTER TABLE nominants DROP CONSTRAINT IF EXISTS nominants_category_id_fkey;
ALTER TABLE nominants DROP COLUMN IF EXISTS category_id;

-- Создаем таблицу для связи many-to-many между номинантами и категориями
CREATE TABLE IF NOT EXISTS nominant_categories (
    nominant_id BIGINT NOT NULL REFERENCES nominants(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (nominant_id, category_id)
);

-- Создаем индексы для производительности
CREATE INDEX IF NOT EXISTS idx_nominant_categories_nominant_id ON nominant_categories(nominant_id);
CREATE INDEX IF NOT EXISTS idx_nominant_categories_category_id ON nominant_categories(category_id);

-- Удаляем старый индекс
DROP INDEX IF EXISTS idx_nominants_category_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Восстанавливаем структуру
DROP TABLE IF EXISTS nominant_categories;

ALTER TABLE nominants ADD COLUMN category_id BIGINT;
ALTER TABLE nominants ADD CONSTRAINT nominants_category_id_fkey 
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_nominants_category_id ON nominants(category_id);

-- +goose StatementEnd

