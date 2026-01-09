-- Миграция 009: Создание таблицы user_clients для маппинга client_id -> user_id
-- Автор: System
-- Дата: 2026-01-09

CREATE TABLE IF NOT EXISTS user_clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_id VARCHAR(255) UNIQUE NOT NULL, -- ID клиента из api-gateway

    -- Информация о клиенте
    client_info JSONB DEFAULT '{}', -- JSON с информацией о клиенте (IP, User-Agent, etc.)

    -- Assigned video service instance
    assigned_instance_id UUID REFERENCES video_service_instances(id) ON DELETE SET NULL,

    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_user_clients_user_id ON user_clients(user_id);
CREATE INDEX IF NOT EXISTS idx_user_clients_client_id ON user_clients(client_id);
CREATE INDEX IF NOT EXISTS idx_user_clients_active ON user_clients(is_active);
CREATE INDEX IF NOT EXISTS idx_user_clients_instance ON user_clients(assigned_instance_id);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_user_clients_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS trigger_update_user_clients_updated_at ON user_clients;
CREATE TRIGGER trigger_update_user_clients_updated_at
    BEFORE UPDATE ON user_clients
    FOR EACH ROW
    EXECUTE FUNCTION update_user_clients_updated_at();

DO $$
BEGIN
    RAISE NOTICE '✅ Таблица user_clients создана';
END $$;
