-- Complete Database Schema
-- This file combines all migrations in order

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    photo_url TEXT,
    bio TEXT,
    role VARCHAR(50) NOT NULL,
    wallet_balance NUMERIC(12,2) DEFAULT 0,
    media_gallery JSONB DEFAULT '[]',
    cover_photo_url TEXT,
    follower_count INTEGER DEFAULT 0,
    instagram_handle VARCHAR(255),
    verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Events table
CREATE TABLE IF NOT EXISTS events (
    id CHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    date TIMESTAMP NOT NULL,
    location VARCHAR(255),
    host_id CHAR(36) REFERENCES users(id) ON DELETE SET NULL,
    category VARCHAR(50),
    images JSONB DEFAULT '[]',
    min_bid NUMERIC(12,2) NOT NULL,
    capacity INTEGER NOT NULL,
    bid_deadline TIMESTAMP NOT NULL,
    status VARCHAR(50) DEFAULT 'upcoming',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Event Attendees table
CREATE TABLE IF NOT EXISTS event_attendees (
    user_id CHAR(36) REFERENCES users(id) ON DELETE CASCADE,
    event_id CHAR(36) REFERENCES events(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, event_id)
);

-- Applications table
CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    event_id CHAR(36) REFERENCES events(id) ON DELETE CASCADE,
    fan_id CHAR(36) REFERENCES users(id) ON DELETE CASCADE,
    bid_amount NUMERIC(12,2) NOT NULL,
    message TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(event_id, fan_id)
);

-- Chat Rooms table
CREATE TABLE IF NOT EXISTS chat_rooms (
    id SERIAL PRIMARY KEY,
    event_id CHAR(36) REFERENCES events(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    chat_room_id INTEGER REFERENCES chat_rooms(id) ON DELETE CASCADE,
    sender_id CHAR(36) REFERENCES users(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add new profile fields to users table if they don't exist
DO $$ 
BEGIN
    -- Add age column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'age') THEN
        ALTER TABLE users ADD COLUMN age INTEGER;
    END IF;

    -- Add gender column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'gender') THEN
        ALTER TABLE users ADD COLUMN gender VARCHAR(50);
    END IF;

    -- Add height column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'height') THEN
        ALTER TABLE users ADD COLUMN height INTEGER;
    END IF;

    -- Add work column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'work') THEN
        ALTER TABLE users ADD COLUMN work TEXT;
    END IF;

    -- Add education column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'education') THEN
        ALTER TABLE users ADD COLUMN education TEXT;
    END IF;

    -- Add education_level column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'education_level') THEN
        ALTER TABLE users ADD COLUMN education_level VARCHAR(50);
    END IF;

    -- Add drinking column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'drinking') THEN
        ALTER TABLE users ADD COLUMN drinking VARCHAR(50);
    END IF;

    -- Add location column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'location') THEN
        ALTER TABLE users ADD COLUMN location TEXT;
    END IF;

    -- Add government_id_url column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'government_id_url') THEN
        ALTER TABLE users ADD COLUMN government_id_url TEXT;
    END IF;

    -- Add verified_at column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'verified_at') THEN
        ALTER TABLE users ADD COLUMN verified_at TIMESTAMP;
    END IF;

    -- Add events_hosted_count column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'events_hosted_count') THEN
        ALTER TABLE users ADD COLUMN events_hosted_count INTEGER NOT NULL DEFAULT 0;
    END IF;

    -- Add mobile_number column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'mobile_number') THEN
        ALTER TABLE users ADD COLUMN mobile_number VARCHAR(20);
    END IF;
END $$; 