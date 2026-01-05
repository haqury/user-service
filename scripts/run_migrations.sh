#!/bin/bash
# Скрипт для управления миграциями базы данных

set -e

DB_NAME="user_service_db"
DB_USER="user_service"
DB_PASSWORD="SecurePass123!"
CONTAINER_NAME="user-service-postgres"

show_help() {
    echo "Использование: $0 [команда]"
    echo ""
    echo "Команды:"
    echo "  up              Применить все миграции"
    echo "  down            Откатить все миграции"
    echo "  status          Показать статус миграций"
    echo "  create <name>   Создать новую миграцию"
    echo "  reset           Сбросить базу данных и применить миграции"
    echo ""
    echo "Примеры:"
    echo "  $0 up           # Применить миграции"
    echo "  $0 status       # Показать статус"
    echo "  $0 create add_new_field   # Создать миграцию"
}

run_migration() {
    local file="$1"
    echo "📄 Применяем: $file"

    if docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME < "scripts/migrations/$file"; then
        echo "✅ Успешно: $file"
        return 0
    else
        echo "❌ Ошибка: $file"
        return 1
    fi
}

command_up() {
    echo "🔄 Применяем все миграции..."

    for migration_file in scripts/migrations/*.sql; do
        if [ -f "$migration_file" ] && [[ "$migration_file" != *"down/"* ]]; then
            local filename=$(basename "$migration_file")
            if ! run_migration "$filename"; then
                echo "❌ Остановка из-за ошибки"
                exit 1
            fi
        fi
    done

    echo "✅ Все миграции применены"
}

command_down() {
    echo "🔄 Откатываем миграции..."

    if [ -f "scripts/migrations/down/001_drop_all_tables.sql" ]; then
        echo "📄 Откатываем все таблицы..."
        docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME < "scripts/migrations/down/001_drop_all_tables.sql"
        echo "✅ Все таблицы удалены"
    else
        echo "⚠️  Файл отката не найден"
    fi
}

command_status() {
    echo "📊 Статус базы данных:"

    # Проверяем подключение
    if ! docker exec $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -c "SELECT 1" &> /dev/null; then
        echo "❌ Не удалось подключиться к базе данных"
        exit 1
    fi

    # Список таблиц
    echo ""
    echo "📋 Таблицы:"
    docker exec $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -c "SELECT
        tablename as \"Таблица\",
        pg_size_pretty(pg_total_relation_size(quote_ident(tablename))) as \"Размер\"
    FROM pg_tables
    WHERE schemaname = 'public'
    ORDER BY tablename;"

    # Количество записей
    echo ""
    echo "📊 Количество записей:"
    docker exec $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -c "SELECT
        'users' as таблица,
        COUNT(*) as записей
    FROM users
    UNION ALL
    SELECT
        'user_services',
        COUNT(*)
    FROM user_services
    UNION ALL
    SELECT
        'routing_logs',
        COUNT(*)
    FROM routing_logs
    UNION ALL
    SELECT
        'auth_tokens',
        COUNT(*)
    FROM auth_tokens;"
}

command_create() {
    local name="$1"
    if [ -z "$name" ]; then
        echo "❌ Укажите имя миграции"
        exit 1
    fi

    local timestamp=$(date +%Y%m%d%H%M%S)
    local filename="scripts/migrations/${timestamp}_${name}.sql"

    echo "-- Миграция: $name" > "$filename"
    echo "-- Автор: $(whoami)" >> "$filename"
    echo "-- Дата: $(date +%Y-%m-%d)" >> "$filename"
    echo "" >> "$filename"
    echo "-- Вверх" >> "$filename"
    echo "" >> "$filename"
    echo "-- Вниз" >> "$filename"
    echo "-- DELETE FROM ..." >> "$filename"

    echo "✅ Создана миграция: $filename"
    echo "ℹ️  Заполните секции 'Вверх' и 'Вниз'"
}

command_reset() {
    echo "🔄 Сбрасываем базу данных..."

    read -p "⚠️  Это удалит все данные. Продолжить? [y/N]: " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Отменено"
        exit 0
    fi

    command_down
    sleep 2
    command_up

    echo "✅ База данных сброшена"
}

# Основная логика
case "$1" in
    "up")
        command_up
        ;;
    "down")
        command_down
        ;;
    "status")
        command_status
        ;;
    "create")
        command_create "$2"
        ;;
    "reset")
        command_reset
        ;;
    *)
        show_help
        ;;
esac
