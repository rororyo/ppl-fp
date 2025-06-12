CREATE TABLE IF NOT EXISTS user_quiz_sessions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
  started_at TIMESTAMPTZ DEFAULT NOW(),
  ended_at TIMESTAMPTZ,
  submitted BOOLEAN,
  auto_submitted BOOLEAN,
  score INTEGER,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE
)