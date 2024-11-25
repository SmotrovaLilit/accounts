CREATE TABLE IF NOT EXISTS accounts
(
    id         UUID PRIMARY KEY,
    owner_id   UUID        NOT NULL,
    currency   VARCHAR(3)  NOT NULL,
    balance    BIGINT      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);