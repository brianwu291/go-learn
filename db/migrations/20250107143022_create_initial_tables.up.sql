-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create ENUMs
CREATE TYPE rate_unit AS ENUM ('second', 'minute', 'hour', 'day');
CREATE TYPE motor_status AS ENUM ('available', 'in_use', 'maintenance', 'retired');
CREATE TYPE motor_category AS ENUM (
  'scooter',
  'sport',
  'cruiser',
  'touring',
  'standard'
);
CREATE TYPE rental_status AS ENUM ('active', 'completed', 'cancelled', 'overdue');
-- Create users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  name VARCHAR(80) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Create rate_plans table
CREATE TABLE rate_plans (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(80) NOT NULL UNIQUE,
  description TEXT,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Create rate_tiers table
CREATE TABLE rate_tiers (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  rate_plan_id UUID NOT NULL REFERENCES rate_plans(id),
  unit_price DECIMAL(10, 2) NOT NULL,
  unit_type rate_unit NOT NULL,
  start_second INTEGER NOT NULL,
  end_second INTEGER,
  -- NULL means unlimited
  priority INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (rate_plan_id, priority)
);
-- Create motors table
CREATE TABLE motors (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  make VARCHAR(50) NOT NULL,
  model VARCHAR(50) NOT NULL,
  year INTEGER NOT NULL,
  category motor_category NOT NULL,
  plate_number VARCHAR(20) NOT NULL UNIQUE,
  rate_plan_id UUID NOT NULL REFERENCES rate_plans(id),
  status motor_status NOT NULL DEFAULT 'available',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Create rentals table
CREATE TABLE rentals (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id),
  motor_id UUID NOT NULL REFERENCES motors(id),
  rate_plan_id UUID NOT NULL REFERENCES rate_plans(id),
  start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  end_time TIMESTAMP WITH TIME ZONE,
  total_amount DECIMAL(10, 2),
  status rental_status NOT NULL DEFAULT 'active',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Create rental_charges table
CREATE TABLE rental_charges (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  rental_id UUID NOT NULL REFERENCES rentals(id),
  rate_tier_id UUID NOT NULL REFERENCES rate_tiers(id),
  unit_price DECIMAL(10, 2) NOT NULL,
  unit_type rate_unit NOT NULL,
  duration_seconds INTEGER NOT NULL,
  amount DECIMAL(10, 2) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Create indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_motors_status ON motors(status);
CREATE INDEX idx_motors_rate_plan_id ON motors(rate_plan_id);
CREATE INDEX idx_rate_tiers_rate_plan_id ON rate_tiers(rate_plan_id);
CREATE INDEX idx_rate_tiers_start_second ON rate_tiers(start_second);
CREATE INDEX idx_rentals_user_id ON rentals(user_id);
CREATE INDEX idx_rentals_motor_id ON rentals(motor_id);
CREATE INDEX idx_rentals_rate_plan_id ON rentals(rate_plan_id);
CREATE INDEX idx_rentals_status ON rentals(status);
CREATE UNIQUE INDEX idx_rentals_active_per_user ON rentals (user_id, status)
WHERE status = 'active';
CREATE UNIQUE INDEX idx_rentals_active_per_motor ON rentals (motor_id, status)
WHERE status = 'active';
CREATE INDEX idx_rental_charges_rental_id ON rental_charges(rental_id);
-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';
-- Create triggers
CREATE TRIGGER update_users_updated_at BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_rate_plans_updated_at BEFORE
UPDATE ON rate_plans FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_rate_tiers_updated_at BEFORE
UPDATE ON rate_tiers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_motors_updated_at BEFORE
UPDATE ON motors FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_rentals_updated_at BEFORE
UPDATE ON rentals FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
