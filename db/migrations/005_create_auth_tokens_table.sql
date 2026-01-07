-- Миграция 005: Создание таблицы auth_tokens
-- Автор: System
-- Дата: 2025-01-05
-- Описание: Токены аутентификации

CREATE TABLE IF NOT EXISTS auth_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    token VARCHAR(500) NOT NULL,
    token_type VARCHAR(50) DEFAULT 'bearer' CHECK (token_type IN ('bearer', 'jwt', 'api_key')),

    -- Срок действия
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    refresh_token VARCHAR(500),
    refresh_expires_at TIMESTAMP WITH TIME ZONE,

    -- Клиентская информация
    client_info JSONB DEFAULT '{}',
    user_agent TEXT,
    ip_address INET,

    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP WITH TIME ZONE,

    UNIQUE(token)
);

-- Индексы для таблицы auth_tokens
CREATE INDEX IF NOT EXISTS idx_auth_tokens_user_id ON auth_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_token ON auth_tokens(token);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_expires_at ON auth_tokens(expires_at);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_auth_tokens_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS trigger_update_auth_tokens_updated_at ON auth_tokens;
CREATE TRIGGER trigger_update_auth_tokens_updated_at
    BEFORE UPDATE ON auth_tokens
    FOR EACH ROW
    EXECUTE FUNCTION update_auth_tokens_updated_at();

-- Логирование
DO $$
BEGIN
    RAISE NOTICE '✅ Таблица auth_tokens создана';
END $$;
