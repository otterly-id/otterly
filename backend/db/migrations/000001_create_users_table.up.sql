-- Set timezone
SET TIMEZONE="Asia/Jakarta";

-- Create user roles enum
CREATE TYPE user_role AS ENUM ('ADMIN', 'USER', 'OWNER');

-- Create users table
CREATE TABLE users (
	id UUID DEFAULT gen_random_uuid() PRIMARY KEY,

	name VARCHAR (50) NOT NULL UNIQUE,
	full_name VARCHAR (100),
	email VARCHAR (254) NOT NULL UNIQUE,
	password_hash BYTEA NOT NULL,
	phone_number VARCHAR (20),
	role user_role NOT NULL,

	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Create a reusable function to update the 'updated_at' timestamp
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to automatically call the function before any update on the users table
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
