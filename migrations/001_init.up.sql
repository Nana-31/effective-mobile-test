CREATE TABLE IF NOT EXISTS subscriptions (
    service_name TEXT,
    price INTEGER NOT NULL,
    user_id UUID NOT NULL,
    id UUID PRIMARY KEY,
    start_date DATE NOT NULL,
    end_date DATE
)