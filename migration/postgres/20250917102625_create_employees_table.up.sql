CREATE TABLE employees (
    id UUID PRIMARY KEY,
    full_name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_employees_email ON employees(email);

INSERT INTO employees (id, full_name, email) VALUES
('019955b8-0981-7c83-9078-c0e021845487', 'Alice Smith', 'alice.smith@example.com'),
('019955b8-0981-7dd9-a087-ace6d5e4906b', 'Bob Johnson', 'bob.johnson@example.com'),
('019955b8-0981-7787-9747-c04068fbbd62', 'Charlie Brown', 'charlie.brown@example.com'),
('019955b8-0981-7324-acd6-9166775bbb71', 'Diana Prince', 'diana.prince@example.com'),
('019955b8-0981-7c16-afc0-8218a5d0131c', 'Eve Adams', 'eve.adams@example.com');