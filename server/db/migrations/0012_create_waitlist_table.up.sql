CREATE TYPE waitlist_status AS ENUM ('waiting', 'notified', 'confirmed', 'expired');

CREATE TABLE waitlists (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date         DATE NOT NULL,
    time_slot    TIME NOT NULL,
    party_size   INT NOT NULL,
    status       waitlist_status NOT NULL DEFAULT 'waiting',
    notified_at  TIMESTAMPTZ,
    expires_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_waitlists_user_id ON waitlists(user_id);
CREATE INDEX idx_waitlists_date ON waitlists(date);
CREATE INDEX idx_waitlists_id ON waitlists(id);
CREATE INDEX idx_waitlists_status ON waitlists(status);
CREATE INDEX idx_waitlists_deleted_at ON waitlists(deleted_at);
CREATE INDEX idx_waitlists_expires_at ON waitlists(expires_at);
CREATE INDEX idx_waitlists_notified_at ON waitlists(notified_at);
CREATE INDEX idx_waitlists_created_at ON waitlists(created_at);
CREATE INDEX idx_waitlists_lookup ON waitlists(date, time_slot, party_size, status);
CREATE INDEX idx_waitlists_check ON waitlists(user_id, date, time_slot);
CREATE INDEX idx_waitlists_position ON waitlists(date, time_slot, status, created_at);