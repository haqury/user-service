-- Миграция 001: Создание расширений
-- Автор: System
-- Дата: 2025-01-05
-- Описание: Создаем необходимые расширения PostgreSQL

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Логирование
DO $$
BEGIN
    RAISE NOTICE '✅ Расширения созданы: uuid-ossp, pgcrypto';
END $$;
