CREATE TABLE investments (
    id UUID PRIMARY KEY,
    loan_id UUID NOT NULL REFERENCES loans(id),
    investor_id VARCHAR(255) NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_investments_loan_id ON investments(loan_id);
