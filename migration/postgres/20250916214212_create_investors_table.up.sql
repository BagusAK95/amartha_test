CREATE TABLE investors (
    id UUID PRIMARY KEY,
    full_name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    balance float8 NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_investors_email ON investors(email);

INSERT INTO investors (id, full_name, email, balance) VALUES
('0199562e-a33a-71d8-a898-937283748233', 'Michael Chen', 'michael.chen@example.com', 1000000),
('0199562f-4b3a-7f8c-8d83-838283848283', 'Sarah Davis', 'sarah.davis@example.com', 2500000),
('0199562f-d33a-71d8-a898-937283748233', 'David Rodriguez', 'david.rodriguez@example.com', 500000);
