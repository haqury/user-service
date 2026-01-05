-- Миграция 001: Инициализация базы данных
-- Автор: System
-- Дата: $(date +%Y-%m-%d)
-- Описание: Создаем базу данных и расширения

-- Создаем базу данных если не существует
SELECT 'CREATE DATABASE user_service_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'user_service_db')\gexec

-- Подключаемся к новой базе данных
\c user_service_db;

-- Создаем расширения
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Логирование
DO $$
BEGIN
    RAISE NOTICE '✅ База данных user_service_db создана с расширениями';
END $$;
