-- Миграция 010: Начальные данные для video_service_instances
-- Автор: System
-- Дата: 2026-01-09

-- Вставляем тестовые инстансы video-service
INSERT INTO video_service_instances (
    name,
    server_url,
    server_port,
    use_ssl,
    stream_endpoint,
    region,
    priority,
    max_capacity,
    allowed_tiers,
    max_bitrate,
    max_resolution,
    codec,
    is_active
) VALUES
-- Premium инстанс для Европы
(
    'premium-eu-1',
    'video-service-premium.example.com',
    9092,
    true,
    '/stream',
    'eu',
    100,
    500,
    ARRAY['premium', 'enterprise'],
    10000,
    1080,
    'h264',
    true
),
-- Базовый инстанс для Европы
(
    'basic-eu-1',
    'video-service.example.com',
    9091,
    true,
    '/stream',
    'eu',
    50,
    1000,
    ARRAY['free', 'basic', 'premium'],
    5000,
    720,
    'h264',
    true
),
-- Premium инстанс для США
(
    'premium-us-1',
    'video-service-premium-us.example.com',
    9092,
    true,
    '/stream',
    'us',
    100,
    500,
    ARRAY['premium', 'enterprise'],
    10000,
    1080,
    'h264',
    true
),
-- Базовый инстанс для США
(
    'basic-us-1',
    'video-service-us.example.com',
    9091,
    true,
    '/stream',
    'us',
    50,
    1000,
    ARRAY['free', 'basic', 'premium'],
    5000,
    720,
    'h264',
    true
),
-- Дефолтный инстанс
(
    'default-1',
    'video-service-default.example.com',
    9091,
    false,
    '/stream',
    'default',
    10,
    2000,
    ARRAY['free', 'basic', 'premium', 'enterprise'],
    3000,
    480,
    'h264',
    true
)
ON CONFLICT (name) DO NOTHING;

DO $$
BEGIN
    RAISE NOTICE '✅ Начальные данные для video_service_instances добавлены';
END $$;
