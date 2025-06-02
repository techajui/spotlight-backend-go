-- Add events_hosted_count column to users table
ALTER TABLE users ADD COLUMN events_hosted_count INTEGER NOT NULL DEFAULT 0; 