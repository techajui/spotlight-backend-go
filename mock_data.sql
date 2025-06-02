-- Users
INSERT INTO users (id, name, email, password, avatar_url, bio, last_login, created_at, updated_at)
VALUES
  (1, 'John Doe', 'john@example.com', '$2a$10$7QJ8QwQwQwQwQwQwQwQwQeQwQwQwQwQwQwQwQwQwQwQwQwQwQwQ', 'https://api.dicebear.com/7.x/avataaars/svg?seed=John', 'Event organizer and tech enthusiast', NOW(), NOW(), NOW()),
  (2, 'Jane Smith', 'jane@example.com', '$2a$10$7QJ8QwQwQwQwQwQwQwQwQeQwQwQwQwQwQwQwQwQwQwQwQwQwQwQ', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Jane', 'Professional photographer', NOW(), NOW(), NOW()),
  (3, 'Mike Johnson', 'mike@example.com', '$2a$10$7QJ8QwQwQwQwQwQwQwQwQeQwQwQwQwQwQwQwQwQwQwQwQwQwQwQ', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Mike', 'Music producer and DJ', NOW(), NOW(), NOW());

-- Events
INSERT INTO events (id, title, description, start_time, end_time, location, creator_id, status, min_bid, current_bid, created_at, updated_at)
VALUES
  (1, 'Tech Conference 2024', 'Annual technology conference featuring the latest innovations', NOW() + INTERVAL '30 days', NOW() + INTERVAL '32 days', 'San Francisco Convention Center', 1, 'active', 1000, 1000, NOW(), NOW()),
  (2, 'Photography Workshop', 'Learn advanced photography techniques from professionals', NOW() + INTERVAL '15 days', NOW() + INTERVAL '16 days', 'Downtown Art Center', 2, 'active', 500, 500, NOW(), NOW()),
  (3, 'Summer Music Festival', 'Three days of non-stop music and entertainment', NOW() + INTERVAL '60 days', NOW() + INTERVAL '63 days', 'Central Park', 3, 'active', 2000, 2000, NOW(), NOW());

-- Chat Rooms
INSERT INTO chat_rooms (id, event_id, status, created_at, updated_at)
VALUES
  (1, 1, 'active', NOW(), NOW()),
  (2, 2, 'active', NOW(), NOW()),
  (3, 3, 'active', NOW(), NOW());

-- Messages
INSERT INTO messages (id, chat_room_id, sender_id, content, created_at, updated_at)
VALUES
  (1, 1, 1, 'Welcome to the chat room! Feel free to ask any questions about the event.', NOW(), NOW()),
  (2, 1, 1, 'Looking forward to this event!', NOW(), NOW()),
  (3, 1, 2, 'Will there be any special guests?', NOW(), NOW()),
  (4, 2, 2, 'Welcome to the chat room! Feel free to ask any questions about the event.', NOW(), NOW()),
  (5, 2, 1, 'Looking forward to this event!', NOW(), NOW()),
  (6, 2, 2, 'Will there be any special guests?', NOW(), NOW()),
  (7, 3, 3, 'Welcome to the chat room! Feel free to ask any questions about the event.', NOW(), NOW()),
  (8, 3, 1, 'Looking forward to this event!', NOW(), NOW()),
  (9, 3, 2, 'Will there be any special guests?', NOW(), NOW()); 