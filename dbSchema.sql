CREATE TYPE patchtype AS ENUM ('edit');

CREATE TABLE IF NOT EXISTS patches (
  time TIMESTAMPTZ NOT NULL,
  patch TEXT NOT NULL,
  convo_id INTEGER,
  user_id INTEGER,
  type PATCHTYPE
);

SELECT create_hypertable('patches', 'time');
