CREATE TYPE reservation_status AS ENUM ('pending', 'confirmed', 'checked_in', 'no_show', 'cancelled');

CREATE TABLE reservations (
    id               BIGSERIAL PRIMARY KEY,
    user_id          BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    table_id         BIGINT NOT NULL,
    date             DATE NOT NULL,
    time_slot        TIME NOT NULL,
    party_size       INT NOT NULL,
    status           reservation_status NOT NULL DEFAULT 'pending',
    special_requests TEXT,
    checked_in_at    TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_reservations_user_id ON reservations(user_id);
CREATE INDEX idx_reservations_table_id ON reservations(table_id);
CREATE INDEX idx_reservations_date ON reservations(date);
CREATE INDEX idx_reservations_table_date ON reservations(table_id, date);
CREATE INDEX idx_reservations_status ON reservations(status);
CREATE INDEX idx_reservations_deleted_at ON reservations(deleted_at);