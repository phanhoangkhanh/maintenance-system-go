CREATE TABLE plan (
    id VARCHAR(36) PRIMARY KEY,
    customer VARCHAR(255),
    number_vehicle VARCHAR(255),
    driver_code VARCHAR(255),
    operator_code VARCHAR(255),
    note TEXT,
    status VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);