INSERT INTO users (
    name,
    email,
    password_hash
)
VALUES (
    'Admin User',
    'admin@crmlite.com',
    'hashed_password'
);

INSERT INTO leads (
    owner_id,
    name,
    email,
    company,
    status
)
SELECT
    id,
    'Acme Lead',
    'lead@acme.com',
    'Acme Inc',
    'NEW'
FROM users
LIMIT 1;