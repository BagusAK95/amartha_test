CREATE TABLE loans (
    id UUID PRIMARY KEY,
    borrower_id UUID NOT NULL,
    principal_amount float8 NOT NULL,
    rate float4 NOT NULL,
    roi float4 NOT NULL,
    state VARCHAR NOT NULL,
    agreement_letter_url VARCHAR,

    validator_employee_id VARCHAR,
    visit_proof_picture_url VARCHAR,
    approval_date TIMESTAMPTZ,

    officer_employee_id VARCHAR,
    signed_agreement_url VARCHAR,
    disbursement_date TIMESTAMPTZ,

    reject_reason VARCHAR,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_loans_borrower_id ON loans(borrower_id);
CREATE INDEX idx_loans_state ON loans(state);
CREATE INDEX idx_loans_approval_date ON loans(approval_date);
CREATE INDEX idx_loans_disbursement_date ON loans(disbursement_date);
