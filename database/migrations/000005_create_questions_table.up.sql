CREATE TABLE IF NOT EXISTS questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    content JSONB NOT NULL,
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE
);