-- +goose Up
CREATE TYPE frequency_type AS ENUM ('hourly', 'daily');

CREATE TABLE subscriptions (
    token UUID NOT NULL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    city TEXT NOT NULL,
    frequency frequency_type NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS subscriptions;
DROP TYPE IF EXISTS frequency_type;
