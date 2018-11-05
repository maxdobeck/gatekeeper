DELETE FROM enrollments;
DELETE FROM members;
DELETE FROM schedules;
DELETE FROM shifts;
DROP TABLE enrollments;
DROP TABLE members CASCADE;
DROP TABLE schedules CASCADE;
DROP TABLE shifts CASCADE;