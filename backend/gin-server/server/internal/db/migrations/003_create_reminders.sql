CREATE TABLE IF NOT EXISTS reminders (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    is_pinned BOOLEAN DEFAULT FALSE,
    user_id TEXT NOT NULL,
    reminder_group_id TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--- if a row in the users table (the parent table) is deleted,
--- all corresponding rows in the current table (the table with the foreign key) that have a matching user_id will also be automatically deleted.
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (reminder_group_id) REFERENCES reminder_groups(id) ON DELETE SET NULL
)