CREATE TYPE user_role AS ENUM ('Admin', 'Driver', 'Manager');
CREATE TYPE license_cat AS ENUM ('A', 'B', 'C');

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role user_role NOT NULL,
    license_category license_cat NOT NULL
);
