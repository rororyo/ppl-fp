CREATE TABLE IF NOT EXISTS user_answers (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
  session_id UUID NOT NULL REFERENCES user_quiz_sessions(id) ON DELETE CASCADE,
  question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  option_id UUID NOT NULL REFERENCES questions_options(id) ON DELETE CASCADE,
  UNIQUE (session_id, question_id)
);
