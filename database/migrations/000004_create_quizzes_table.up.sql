CREATE TABLE IF NOT EXISTS quizzes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quiz_name TEXT NOT NULL,
    time_limit INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE
);
