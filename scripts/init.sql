CREATE USER tasker;
CREATE DATABASE tasks;
GRANT ALL PRIVILEGES ON DATABASE tasks TO tasker;
CREATE TABLE tasks (id serial, name char(124));
