-- Миграция 006: Начальные данные
-- Автор: System
-- Дата: 2025-01-05
-- Описание: Добавляем тестового пользователя и сервисы

-- Тестовый пользователь
INSERT INTO users (username, email, phone, password_hash) VALUES
('testuser', 'test@example.com', '+1234567890', crypt('password123', gen_salt('bf', 10)))
ON CONFLICT (username) DO NOTHING;

-- Получаем ID созданного пользователя и добавляем сервисы
DO $$
DECLARE
    user_uuid UUID;
BEGIN
    -- Получаем ID пользователя
    SELECT id INTO user_uuid FROM users WHERE username = 'testuser';

    IF user_uuid IS NOT NULL THEN
        -- Добавляем сервисы для перенаправления (3 кнопки)
        INSERT INTO user_services (user_id, service_name, service_type, base_url, port, priority, enabled_buttons) VALUES
        -- Кнопка 1: Основной стриминг
        (user_uuid, 'primary_stream', 'streaming', 'http://streaming-service-1.internal', 8081, 1, '["button1"]'),

        -- Кнопка 2: Резервный стриминг
        (user_uuid, 'backup_stream', 'streaming', 'http://streaming-service-2.internal', 8082, 2, '["button2"]'),

        -- Кнопка 3: Запись и аналитика
        (user_uuid, 'recording_analytics', 'recording', 'http://recording-service.internal', 8083, 3, '["button3"]')
        ON CONFLICT (user_id, service_name) DO NOTHING;

        RAISE NOTICE '✅ Добавлены сервисы для пользователя testuser';
    ELSE
        RAISE NOTICE 'ℹ️  Пользователь testuser не найден (возможно уже существует)';
    END IF;
END $$;

-- Проверяем созданные данные
SELECT '=== ПОЛЬЗОВАТЕЛИ ===' as info;
SELECT id, username, email, status, created_at FROM users;

SELECT '=== СЕРВИСЫ ДЛЯ AGW ===' as info;
SELECT
    us.service_name as "Сервис",
    us.base_url as "URL",
    us.port as "Порт",
    us.service_type as "Тип",
    us.priority as "Приоритет",
    us.enabled_buttons as "Кнопки"
FROM user_services us
JOIN users u ON us.user_id = u.id
WHERE u.username = 'testuser'
ORDER BY us.priority;
