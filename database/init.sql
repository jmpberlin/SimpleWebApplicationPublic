-- Initial database schema
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY,
    message TEXT NOT NULL
);

-- Seed data for the three buttons
INSERT INTO messages (id, message) VALUES
    (1, 'Hello from the database!'),
    (2, 'Goodbye from the database!'),
    (3, 'Knock knock! Who is there? The database!')
ON CONFLICT (id) DO NOTHING;

-- Log completion
SELECT 'Database initialization complete!' as status;
