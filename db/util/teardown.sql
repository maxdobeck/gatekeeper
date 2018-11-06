DELETE FROM enrollments;
DELETE FROM members;
DELETE FROM schedules;
DELETE FROM shifts;
DROP TABLE enrollments;
DROP TABLE shifts CASCADE;
DROP TABLE schedules CASCADE;
DROP TABLE members CASCADE;