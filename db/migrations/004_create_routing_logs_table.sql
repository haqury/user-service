-- Миграция 004: Создание таблицы routing_logs
-- Автор: System
-- Дата: 2025-01-05
-- Описание: Логи перенаправлений AGW

CREATE TABLE IF NOT EXISTS routing_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_id UUID REFERENCES user_services(id) ON DELETE SET NULL,

    -- Данные запроса
    request_id VARCHAR(100) NOT NULL,
    source_ip INET,
    user_agent TEXT,

    -- Маршрутизация
    target_service VARCHAR(100),
    target_url VARCHAR(500),
    target_port INTEGER,

    -- Данные
    data_size_bytes INTEGER,
    data_type VARCHAR(50),

    -- Статус
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'success', 'failed', 'retry')),
    status_code INTEGER,
    error_message TEXT,

    -- Временные метки
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_ms INTEGER,

    -- Дополнительные данные
    metadata JSONB DEFAULT '{}'
);

-- Индексы для таблицы routing_logs
CREATE INDEX IF NOT EXISTS idx_routing_logs_user_id ON routing_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_routing_logs_service_id ON routing_logs(service_id);
CREATE INDEX IF NOT EXISTS idx_routing_logs_request_id ON routing_logs(request_id);
CREATE INDEX IF NOT EXISTS idx_routing_logs_status ON routing_logs(status);
CREATE INDEX IF NOT EXISTS idx_routing_logs_started_at ON routing_logs(started_at);

-- Логирование
DO $$
BEGIN
    RAISE NOTICE '✅ Таблица routing_logs создана (AGW logs)';
END $$;
