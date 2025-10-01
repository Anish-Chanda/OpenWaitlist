-- ENUMS
CREATE TYPE auth_provider AS ENUM ('local', 'google');



-- TABLES
CREATE TABLE users (
  id             SERIAL PRIMARY KEY,
  email          TEXT UNIQUE NOT NULL,
  auth_provider  auth_provider NOT NULL,
  password_hash  TEXT,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  display_name   TEXT
);

CREATE TABLE waitlists (
  id                     SERIAL PRIMARY KEY,
  slug                   TEXT UNIQUE NOT NULL,
  name                   TEXT NOT NULL,
  owner_user_id          INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  is_public              BOOLEAN NOT NULL DEFAULT TRUE,
  show_vendor_branding   BOOLEAN NOT NULL DEFAULT FALSE, -- false = remove badge and other open-waitlist mentions
  created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
  archived_at            TIMESTAMPTZ
);