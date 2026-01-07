#!/bin/bash

# Проверка структуры миграций
echo "🔍 Checking migration structure..."

# Проверяем существование директории
if [ ! -d "db/migrations" ]; then
    echo "❌ Directory db/migrations not found"
    exit 1
fi

# Считаем миграции
migration_count=$(ls -1 db/migrations/*.sql 2>/dev/null | wc -l)
echo "✅ Migration structure is valid"
echo "📊 Total migrations: $migration_count"

echo "\n📋 Available migrations (in execution order):"
for file in $(ls db/migrations/*.sql | sort); do
    if [ -f "$file" ]; then
        echo "  📄 $(basename $file)"
    fi
done

# Проверяем формат имен файлов
echo "\n🔎 Checking file naming format..."
bad_files=0
for file in db/migrations/*.sql; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        if [[ ! "$filename" =~ ^[0-9]{14}_.+\.sql$ ]] && [[ ! "$filename" =~ ^[0-9]{3}_.+\.sql$ ]]; then
            echo "  ⚠  Non-standard name: $filename"
            ((bad_files++))
        fi
    fi
done

if [ $bad_files -eq 0 ]; then
    echo "✅ All migration files have proper naming"
else
    echo "⚠  Some files have non-standard names"
fi
