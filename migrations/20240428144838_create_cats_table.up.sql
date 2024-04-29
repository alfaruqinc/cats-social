CREATE TABLE IF NOT EXISTS cats (
    id uuid PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULl,
    name VARCHAR(30) NOT NULL,
    race VARCHAR(24) NOT NULL,
    sex VARCHAR(6) NOT NULL,
    age_in_month INT NOT NULL,
    description VARCHAR(200) NOT NULL,
    image_urls TEXT[] NOT NULL
);
