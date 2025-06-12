CREATE TABLE IF NOT EXISTS questions_options (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
  option TEXT NOT NULL,
  question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE
)