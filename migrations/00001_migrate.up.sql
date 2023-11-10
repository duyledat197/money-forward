CREATE TYPE role_type AS ENUM ('SUPER_ADMIN', 'ADMIN', 'USER');

CREATE TABLE IF NOT EXISTS users (
  ID BIGINT PRIMARY,
  user_name TEXT UNIQUE,
  password TEXT NOT NULL,
  "name" TEXT,
  created_by TEXT,
  role role_type NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT NOW(),
  FOREIGN KEY ("created_by") REFERENCES "users"("id") ON
  DELETE
    CASCADE
);

CREATE TABLE IF NOT EXISTS accounts (
  ID BIGINT PRIMARY,
  user_id TEXT,
  NAME TEXT,
  balance BIGINT,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT NOW(),
  FOREIGN KEY ("user_id") REFERENCES "users"("id") ON
  DELETE
    CASCADE
);