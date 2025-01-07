-- Drop triggers
DROP TRIGGER IF EXISTS update_rentals_updated_at ON rentals;
DROP TRIGGER IF EXISTS update_motors_updated_at ON motors;
DROP TRIGGER IF EXISTS update_rate_tiers_updated_at ON rate_tiers;
DROP TRIGGER IF EXISTS update_rate_plans_updated_at ON rate_plans;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();
-- Drop indexes
DROP INDEX IF EXISTS idx_rental_charges_rental_id;
DROP INDEX IF EXISTS idx_rentals_status;
DROP INDEX IF EXISTS idx_rentals_rate_plan_id;
DROP INDEX IF EXISTS idx_rentals_motor_id;
DROP INDEX IF EXISTS idx_rentals_user_id;
DROP INDEX IF EXISTS idx_rate_tiers_start_second;
DROP INDEX IF EXISTS idx_rate_tiers_rate_plan_id;
DROP INDEX IF EXISTS idx_motors_rate_plan_id;
DROP INDEX IF EXISTS idx_motors_status;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_rentals_active_per_user;
DROP INDEX IF EXISTS idx_rentals_active_per_motor;
-- Drop tables (in correct order due to dependencies)
DROP TABLE IF EXISTS rental_charges;
DROP TABLE IF EXISTS rentals;
DROP TABLE IF EXISTS motors;
DROP TABLE IF EXISTS rate_tiers;
DROP TABLE IF EXISTS rate_plans;
DROP TABLE IF EXISTS users;
-- Drop ENUMs
DROP TYPE IF EXISTS rental_status;
DROP TYPE IF EXISTS motor_category;
DROP TYPE IF EXISTS motor_status;
DROP TYPE IF EXISTS rate_unit;
-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
