-- Add new profile fields to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS age INTEGER,
ADD COLUMN IF NOT EXISTS gender VARCHAR(50),
ADD COLUMN IF NOT EXISTS height INTEGER,
ADD COLUMN IF NOT EXISTS work TEXT,
ADD COLUMN IF NOT EXISTS education TEXT,
ADD COLUMN IF NOT EXISTS education_level VARCHAR(50),
ADD COLUMN IF NOT EXISTS drinking VARCHAR(50),
ADD COLUMN IF NOT EXISTS location TEXT,
ADD COLUMN IF NOT EXISTS government_id_url TEXT,
ADD COLUMN IF NOT EXISTS verified_at TIMESTAMP;

-- Update mock data with new fields
UPDATE users
SET 
    age = 25,
    gender = 'male',
    height = 180,
    work = 'Software Engineer',
    education = 'Bachelor''s in Computer Science',
    education_level = 'bachelors',
    drinking = 'social',
    location = 'San Francisco, CA'
WHERE id = '1';

UPDATE users
SET 
    age = 28,
    gender = 'female',
    height = 165,
    work = 'Professional Photographer',
    education = 'Master''s in Fine Arts',
    education_level = 'masters',
    drinking = 'rarely',
    location = 'New York, NY'
WHERE id = '2';

UPDATE users
SET 
    age = 30,
    gender = 'male',
    height = 175,
    work = 'Music Producer',
    education = 'Bachelor''s in Music Production',
    education_level = 'bachelors',
    drinking = 'yes',
    location = 'Los Angeles, CA'
WHERE id = '3'; 