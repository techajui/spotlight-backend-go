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