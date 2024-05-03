CREATE TABLE IF NOT EXISTS
    users
    (
        id UUID PRIMARY KEY,
        email VARCHAR(100) NOT NULL,
        name VARCHAR(100) NOT NULL,
        password TEXT NOT NULL,
        createdAt TIMESTAMPTZ DEFAULT now(),
            CONSTRAINT
                unique_email
                    UNIQUE(email)

    );
