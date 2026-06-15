CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
);

CREATE TABLE IF NOT EXISTS services_name (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    service_id INTEGER NOT NULL REFERENCES services_name(id),
    price INTEGER NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    start_date DATE NOT NULL DEFAULT CURRENT_DATE,
    end_date DATE
);
