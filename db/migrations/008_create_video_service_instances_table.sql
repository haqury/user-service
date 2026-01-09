-- Миграция 008: Создание таблицы video_service_instances
-- Автор: System
-- Дата: 2026-01-09

CREATE TABLE IF NOT EXISTS video_service_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL, -- название инстанса (например "premium-eu", "basic-us")
    server_url VARCHAR(255) NOT NULL, -- адрес video-service
    server_port INT NOT NULL, -- порт video-service
    use_ssl BOOLEAN DEFAULT true, -- использовать SSL
    stream_endpoint VARCHAR(255) DEFAULT '/stream', -- endpoint для стрима
    region VARCHAR(50) DEFAULT 'default', -- регион инстанса

    -- Параметры для роутинга
    priority INT DEFAULT 0, -- приоритет (больше = выше приоритет)
    max_capacity INT DEFAULT 1000, -- максимальное количество одновременных стримов
    current_load INT DEFAULT 0, -- текущая нагрузка
    health_status VARCHAR(20) DEFAULT 'healthy' CHECK (health_status IN ('healthy', 'degraded', 'unhealthy')),

    -- Подходит для тарифов
    allowed_tiers TEXT[] DEFAULT ARRAY['free', 'basic', 'premium', 'enterprise'],

    -- Лимиты для этого инстанса
    max_bitrate INT DEFAULT 5000, -- максимальный битрейт
    max_resolution INT DEFAULT 1080, -- максимальное разрешение
    codec VARCHAR(50) DEFAULT 'h264', -- поддерживаемый кодек

    -- Метаданные
    metadata JSONB DEFAULT '{}',

    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_health_check TIMESTAMP WITH TIME ZONE
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_video_service_instances_region ON video_service_instances(region);
CREATE INDEX IF NOT EXISTS idx_video_service_instances_health ON video_service_instances(health_status);
CREATE INDEX IF NOT EXISTS idx_video_service_instances_active ON video_service_instances(is_active);
CREATE INDEX IF NOT EXISTS idx_video_service_instances_tiers ON video_service_instances USING GIN(allowed_tiers);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_video_service_instances_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS trigger_update_video_service_instances_updated_at ON video_service_instances;
CREATE TRIGGER trigger_update_video_service_instances_updated_at
    BEFORE UPDATE ON video_service_instances
    FOR EACH ROW
    EXECUTE FUNCTION update_video_service_instances_updated_at();

DO $$
BEGIN
    RAISE NOTICE '✅ Таблица video_service_instances создана';
END $$;
