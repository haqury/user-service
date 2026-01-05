-- Миграция 004: Создание остальных таблиц
-- Автор: System
-- Дата: $(date +%Y-%m-%d)

\c user_service_db;

-- Таблица routing_logs
CREATE TABLE IF NOT EXISTS routing_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_id UUID REFERENCES user_services(id) ON DELETE SET NULL,

    request_id VARCHAR(100) NOT NULL,
    source_ip INET,
    user_agent TEXT,

    target_service VARCHAR(100),
    target_url VARCHAR(500),
    target_port INTEGER,

    data_size_bytes INTEGER,
    data_type VARCHAR(50),

    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'success', 'failed', 'retry')),
    status_code INTEGER,
    error_message TEXT,

    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_ms INTEGER,

    metadata JSONB DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_routing_logs_user_id ON routing_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_routing_logs_service_id ON routing_logs(service_id);
CREATE INDEX IF NOT EXISTS idx_routing_logs_request_id ON routing_logs(request_id);
CREATE INDEX IF NOT EXISTS idx_routing_logs_status ON routing_logs(status);
CREATE INDEX IF NOT EXISTS idx_routing_logs_started_at ON routing_logs(started_at);

-- Таблица auth_tokens
CREATE TABLE IF NOT EXISTS auth_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    token VARCHAR(500) NOT NULL,
    token_type VARCHAR(50) DEFAULT 'bearer' CHECK (token_type IN ('bearer', 'jwt', 'api_key')),

    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    refresh_token VARCHAR(500),
    refresh_expires_at TIMESTAMP WITH TIME ZONE,

    client_info JSONB DEFAULT '{}',
    user_agent TEXT,
    ip_address INET,

    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP WITH TIME ZONE,

    UNIQUE(token)
);

CREATE INDEX IF NOT EXISTS idx_auth_tokens_user_id ON auth_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_token ON auth_tokens(token);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_expires_at ON auth_tokens(expires_at);

-- Триггер
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

DO $$
BEGIN
    RAISE NOTICE '✅ Все таблицы созданы';
END $$;
