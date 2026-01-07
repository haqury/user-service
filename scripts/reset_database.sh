#!/bin/bash

# Очистка базы данных (только для разработки!)
echo "⚠️  WARNING: This will DROP ALL TABLES in the database!"
read -p "Are you sure? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "❌ Operation cancelled"
    exit 1
fi

# Загружаем конфигурацию
echo "🔧 Loading configuration..."
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-user_service}

# Подключаемся к БД и удаляем все таблицы
echo "🗑️  Dropping all tables..."
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" << EOF
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO $DB_USER;
GRANT ALL ON SCHEMA public TO public;
EOF

if [ $? -eq 0 ]; then
    echo "✅ Database reset successfully"
    echo "📋 Next steps:"
    echo "   1. Run migrations: make migrate"
    echo "   2. Start service: make dev"
else
    echo "❌ Failed to reset database"
    exit 1
fi