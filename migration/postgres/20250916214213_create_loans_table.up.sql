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

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
