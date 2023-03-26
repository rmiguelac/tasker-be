CREATE USER tasker;
CREATE DATABASE tasks;
GRANT ALL PRIVILEGES ON DATABASE tasks TO tasker;
ALTER USER tasker WITH PASSWORD 'taskerPWD22';
\c tasks tasker;
CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(124),
    createdAt TIMESTAMP DEFAULT NOW(),
    lastUpdated TIMESTAMP,
    finishedAt TIMESTAMP,
    done BOOLEAN DEFAULT FALSE,
    description TEXT
);
CREATE TABLE IF not EXISTS tags (
    id BIGSERIAL PRIMARY KEY,
    tag VARCHAR(124) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS task_tags (
    task_id INT REFERENCES tasks (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tag_id INT REFERENCES tags (id) ON UPDATE CASCADE
);