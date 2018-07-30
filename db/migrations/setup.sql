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
  owner_id uuid REFERENCES members(id)
);

CREATE TABLE IF NOT EXISTS enrollments(
	schedule_id uuid REFERENCES schedules,
	member_id uuid REFERENCES members,
	admin BOOLEAN CONSTRAINT admin_powers NOT NULL,
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4()
);

