-- Add mobile_number column to users table
ALTER TABLE users ADD COLUMN mobile_number VARCHAR(20);

-- Update existing records with default value
UPDATE users SET mobile_number = '123' WHERE mobile_number IS NULL; 