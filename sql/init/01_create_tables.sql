CREATE TABLE IF NOT EXISTS loan (
    id                bigserial PRIMARY KEY,
    user_id           bigint NOT NULL,
    status            smallint NOT NULL,
    loan_amount       numeric NOT NULL,
    total_outstanding numeric NOT NULL,
    total_paid        numeric NOT NULL,
    total_week        smallint NOT NULL,
    create_time       timestamptz NOT NULL,
    update_time       timestamptz
);

CREATE TABLE IF NOT EXISTS billing (
    id           bigserial PRIMARY KEY,
    loan_id      bigint NOT NULL REFERENCES loan (id),
    status       smallint NOT NULL,
    amount       numeric NOT NULL,
    due_date     timestamptz NOT NULL,
    payment_time timestamptz
);

CREATE INDEX IF NOT EXISTS idx_loan_user_id_create_time ON loan (user_id, create_time);

CREATE INDEX IF NOT EXISTS idx_billing_loan_id ON billing (loan_id);
CREATE INDEX IF NOT EXISTS idx_billing_due_date ON billing (due_date,status);
