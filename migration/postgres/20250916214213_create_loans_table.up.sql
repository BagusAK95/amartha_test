CREATE TABLE loans (
    id UUID PRIMARY KEY,
    borrower_id VARCHAR(255) NOT NULL,
    principal_amount NUMERIC(15, 2) NOT NULL,
    rate NUMERIC(5, 2) NOT NULL,
    roi NUMERIC(5, 2) NOT NULL,
    state VARCHAR(50) NOT NULL,
    agreement_letter_url VARCHAR(255),

    field_validator_employee_id VARCHAR(255),
    visit_proof_picture_url VARCHAR(255),
    approval_date TIMESTAMPTZ,

    field_officer_employee_id VARCHAR(255),
    signed_agreement_url VARCHAR(255),
    disbursement_date TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
