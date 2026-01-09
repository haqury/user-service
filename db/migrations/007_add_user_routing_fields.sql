-- Миграция 007: Добавление полей для роутинга и ролей пользователей
-- Автор: System
-- Дата: 2026-01-09

-- Добавляем новые поля в таблицу users
ALTER TABLE users
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true,
ADD COLUMN IF NOT EXISTS roles TEXT[] DEFAULT ARRAY['user'],
ADD COLUMN IF NOT EXISTS subscription_tier VARCHAR(50) DEFAULT 'free' CHECK (subscription_tier IN ('free', 'basic', 'premium', 'enterprise')),
ADD COLUMN IF NOT EXISTS region VARCHAR(50) DEFAULT 'default';

-- Индексы для новых полей
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_subscription_tier ON users(subscription_tier);
CREATE INDEX IF NOT EXISTS idx_users_region ON users(region);
CREATE INDEX IF NOT EXISTS idx_users_roles ON users USING GIN(roles);

DO $$
BEGIN
    RAISE NOTICE '✅ Добавлены поля для роутинга пользователей';
END $$;
