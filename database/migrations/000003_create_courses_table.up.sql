CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    course_name TEXT NOT NULL,
    content JSONB NOT NULL,
    grade_level INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE
);