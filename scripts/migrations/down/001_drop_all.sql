-- Откат всех таблиц
DROP TABLE IF EXISTS auth_tokens CASCADE;
DROP TABLE IF EXISTS routing_logs CASCADE;
DROP TABLE IF EXISTS user_services CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP FUNCTION IF EXISTS update_users_updated_at CASCADE;
DROP FUNCTION IF EXISTS update_user_services_updated_at CASCADE;
DROP FUNCTION IF EXISTS update_auth_tokens_updated_at CASCADE;
