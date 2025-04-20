## Relationship

- `User` <--1:N--> `ReminderGroup`
- `User` <--1:N--> `Reminder`
- `Reminder` <--N:1--> `User` (already covered by User <--1:N--> Reminder)
- `ReminderGroup` <--1:N--> `Reminder` (optional, a Reminder can exist without a group)

## SQL Databases (SQLite, PostgreSQL)

- Foreign Keys: Foreign keys are the standard way to represent relationships in relational databases.
- In the `reminder_groups` table, we might consider adding a `user_id` as a foreign key to link it back to the users table (to enforce "user creates reminder group").
- In the `reminders` table, we must have a `user_id` as a foreign key referencing the `users` table.
- In the `reminders` table, we can have an optional `reminder_group_id` as a foreign key referencing the reminder_groups table. Making it optional is key to allow reminders without groups (nullable foreign key).


## NoSQL Databases (MongoDB)

Embedding vs. Referencing: MongoDB is flexible. We have options:

- Referencing (Normalization): Similar to foreign keys in SQL, we can store IDs of related documents.
  - `reminder_groups` can have a `user_id` field storing the `_id` of the user who created it.
  - `reminders` must have a `user_id` field referencing the `_id` of the user.
  - `reminders` can have a `reminder_group_id` field, which can be `null` or absent if the reminder isn't in a group, referencing `reminder_groups._id`.
- Embedding (Denormalization): We could embed related data within a document. However, for these relationships (especially User to Reminder/ReminderGroup which are 1-to-N), referencing is generally better to avoid data duplication and maintain consistency when users or groups are updated independently. Embedding would become complex and less efficient for queries and updates.

Relationship Breakdown:

1. User creates Reminder Groups (One-to-Many):
   1. A `User` can create multiple `ReminderGroups`.
   2. A `ReminderGroup` is created by only one `User`.

2. User creates Reminders (One-to-Many):
   1. A `User` can create multiple `Reminders`.
   2. A `Reminder` is always created by and belongs to one `User`.

3. Reminder can belong to a ReminderGroup (Optional Many-to-One):
   1. A `Reminder` can be associated with at most one `ReminderGroup`. It's optional; a reminder doesn't have to belong to a group.
   2. A `ReminderGroup` can have multiple `Reminders` associated with it.
