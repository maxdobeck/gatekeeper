CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS members(
  name text CONSTRAINT name_present NOT NULL,
  email text CONSTRAINT email_present NOT NULL UNIQUE,
  password text CONSTRAINT password_present NOT NULL,
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4()
);

CREATE TABLE IF NOT EXISTS schedules(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  title text CONSTRAINT title_present NOT NULL,
  owner_id uuid REFERENCES members(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS enrollments(
	schedule_id uuid REFERENCES schedules ON DELETE CASCADE,
	member_id uuid REFERENCES members ON DELETE CASCADE,
	admin BOOLEAN CONSTRAINT admin_powers NOT NULL,
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4()
);

CREATE TABLE IF NOT EXISTS shifts(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  schedule_id uuid REFERENCES schedules(id) ON DELETE CASCADE,
  title text CONSTRAINT title_present NOT NULL,
  min_enrollees integer DEFAULT 0,
  start_time text CONSTRAINT start_time_present NOT NULL,
  end_time text CONSTRAINT end_time_present NOT NULL,
  stop_date text DEFAULT "2099-01-01",
  sun boolean DEFAULT FALSE,
  mon boolean DEFAULT FALSE,
  tue boolean DEFAULT FALSE,
  wed boolean DEFAULT FALSE,
  thu boolean DEFAULT FALSE,
  fri boolean DEFAULT FALSE,
  sat boolean DEFAULT FALSE
);