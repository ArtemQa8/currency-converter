-- +goose Up
CREATE TABLE rate_history (
    id BIGSERIAL PRIMARY KEY,
    base_currency VARCHAR(3) NOT NULL,
    quote_currency VARCHAR(3) NOT NULL,
    value NUMERIC(15, 6) NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    source VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ind_rate_history_pair_date ON rate_history (base_currency, quote_currency, date DESC);

-- +goose Down
DROP TABLE rate_history;