CREATE TABLE borrowers (
    id UUID PRIMARY KEY,
    full_name VARCHAR NOT NULL,
    id_card_number VARCHAR UNIQUE NOT NULL,
    address TEXT NOT NULL,
    phone_number VARCHAR UNIQUE NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_borrowers_status ON borrowers(status);
CREATE INDEX idx_borrowers_email ON borrowers(email);
CREATE INDEX idx_borrowers_phone_number ON borrowers(phone_number);
CREATE INDEX idx_borrowers_id_card_number ON borrowers(id_card_number);

INSERT INTO borrowers (id, full_name, id_card_number, address, phone_number, email, status) VALUES
('01995594-af29-71a8-a979-62cdd1419267', 'John Doe', '2834823546872354', '123 Main St', '08123456789', 'john.doe@example.com', 'active'),
('01995594-af29-7a92-9a1b-106b931dfda9', 'Jane Smith', '3201234567890001', '456 Oak Ave', '08123456780', 'jane.smith@example.com', 'active'),
('01995594-af29-7bee-8028-59ff34d746fd', 'Peter Jones', '3201234567890002', '789 Pine Ln', '08123456781', 'peter.jones@example.com', 'inactive'),
('01995594-af29-709c-bf5f-96cb28d01412', 'Mary Johnson', '3201234567890003', '101 Maple Dr', '08123456782', 'mary.johnson@example.com', 'active'),
('01995594-af29-7d57-a3fc-8553b4421c67', 'David Williams', '3201234567890004', '212 Birch Ct', '08123456783', 'david.williams@example.com', 'active');
