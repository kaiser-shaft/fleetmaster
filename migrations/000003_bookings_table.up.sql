CREATE TYPE booking_status AS ENUM('Pending', 'Active', 'Completed', 'Cancelled');

CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    vehicle_id BIGINT NOT NULL REFERENCES vehicles(id),
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    purpose VARCHAR(255),
    status booking_status NOT NULL DEFAULT 'Pending',

    CONSTRAINT duration_check CHECK (end_time >= start_time + INTERVAL '30 minutes')
);

CREATE UNIQUE INDEX idx_user_active_booking ON bookings(user_id)
WHERE (status IN ('Pending', 'Active'));

CREATE INDEX idx_booking_vehicle_id ON bookings(vehicle_id);
