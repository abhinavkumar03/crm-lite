CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(255) NOT NULL,

    email VARCHAR(255) UNIQUE NOT NULL,

    password_hash TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE leads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    owner_id UUID NOT NULL REFERENCES users(id),

    name VARCHAR(255) NOT NULL,

    email VARCHAR(255),

    phone VARCHAR(50),

    company VARCHAR(255),

    status VARCHAR(50) NOT NULL DEFAULT 'NEW',

    notes TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contacts
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    first_name VARCHAR(100) NOT NULL,

    last_name VARCHAR(100),

    email VARCHAR(255),

    phone VARCHAR(20),

    company VARCHAR(150),

    job_title VARCHAR(150),

    notes TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tasks
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    owner_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    lead_id UUID
        REFERENCES leads(id)
        ON DELETE SET NULL,

    contact_id UUID
        REFERENCES contacts(id)
        ON DELETE SET NULL,

    title VARCHAR(200) NOT NULL,

    description TEXT,

    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',

    due_date TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE activity_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    lead_id UUID NOT NULL REFERENCES leads(id),

    activity_type VARCHAR(100) NOT NULL,

    message TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE notes (

    id UUID PRIMARY KEY,

    entity_type VARCHAR(30) NOT NULL,

    entity_id UUID NOT NULL,

    note TEXT NOT NULL,

    created_by UUID NOT NULL,

    updated_by UUID,

    created_at TIMESTAMP NOT NULL,

    updated_at TIMESTAMP NOT NULL

);