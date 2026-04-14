CREATE TYPE vehicle_status AS ENUM('Available', 'In_Use', 'Maintenance', 'Retired');

CREATE TABLE IF NOT EXISTS vehicles (
    id BIGSERIAL PRIMARY KEY,
    brand VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    plate_number VARCHAR(15) NOT NULL UNIQUE,
    status vehicle_status NOT NULL DEFAULT 'Available',
    mileage INTEGER DEFAULT 0 CHECK (mileage >= 0),
    last_service_mileage INTEGER DEFAULT 0 CHECK (last_service_mileage >= 0),

    CONSTRAINT check_mileage_logic CHECK (mileage >= last_service_mileage)
);
