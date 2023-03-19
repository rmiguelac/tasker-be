CREATE USER tasker;
CREATE DATABASE tasks;
GRANT ALL PRIVILEGES ON DATABASE tasks TO tasker;
ALTER USER tasker WITH PASSWORD 'taskerPWD22';
\c tasks tasker;
CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(124),
    createdAt TIMESTAMP DEFAULT NOW(),
    lastUpdated TIMESTAMP DEFAULT NOW(),
    finishedAt TIMESTAMP,
    done BOOLEAN DEFAULT FALSE,
    description TEXT
);
