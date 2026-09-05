CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'Новая',
    user_uid TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);