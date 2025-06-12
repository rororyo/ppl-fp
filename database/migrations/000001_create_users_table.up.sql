CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('admin', 'user');

END IF;

END $$;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    grade_level INTEGER NOT NULL,
    role user_role DEFAULT 'user',
    avatar_url TEXT,
    birth_date DATE,
    token TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);