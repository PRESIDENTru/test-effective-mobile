CREATE TABLE IF NOT EXISTS services_name (
                                             id SERIAL PRIMARY KEY,
                                             name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS services (
                                        id SERIAL PRIMARY KEY,
                                        service_id INTEGER NOT NULL REFERENCES services_name(id),
    price INTEGER NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
    );