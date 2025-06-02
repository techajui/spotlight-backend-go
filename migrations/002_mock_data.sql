-- Insert sample users
INSERT INTO users (name, email, password, avatar_url, bio, last_login) VALUES
('John Doe', 'john@example.com', '$2a$10$abcdefghijklmnopqrstuvwxyz', 'https://example.com/avatar1.jpg', 'Event enthusiast', NOW()),
('Jane Smith', 'jane@example.com', '$2a$10$abcdefghijklmnopqrstuvwxyz', 'https://example.com/avatar2.jpg', 'Sports lover', NOW()),
('Alice Johnson', 'alice@example.com', '$2a$10$abcdefghijklmnopqrstuvwxyz', 'https://example.com/avatar3.jpg', 'Fitness guru', NOW());

-- Insert sample events
INSERT INTO events (title, description, start_time, end_time, location, creator_id, status, min_bid, current_bid) VALUES
('Summer Soccer Tournament', 'Annual soccer tournament for all ages', NOW() + INTERVAL '1 day', NOW() + INTERVAL '2 days', 'Central Park', 1, 'active', 100.00, 150.00),
('Basketball Championship', 'City-wide basketball championship', NOW() + INTERVAL '3 days', NOW() + INTERVAL '4 days', 'Sports Complex', 2, 'active', 200.00, 250.00),
('Tennis Open', 'Open tennis tournament for amateurs', NOW() + INTERVAL '5 days', NOW() + INTERVAL '6 days', 'Tennis Club', 3, 'active', 50.00, 75.00);

-- Insert sample chat rooms
INSERT INTO chat_rooms (event_id, status) VALUES
(1, 'active'),
(2, 'active'),
(3, 'active');

-- Insert sample messages
INSERT INTO messages (chat_room_id, sender_id, content, read_at) VALUES
(1, 1, 'Welcome to the Summer Soccer Tournament chat!', NOW()),
(1, 2, 'Looking forward to the event!', NOW()),
(2, 2, 'Basketball Championship details here.', NOW()),
(2, 3, 'Count me in!', NOW()),
(3, 3, 'Tennis Open registration is open.', NOW()),
(3, 1, 'I will join!', NOW()); 