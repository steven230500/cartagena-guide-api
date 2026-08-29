CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    type TEXT NOT NULL,
    price TEXT NOT NULL,
    price_amount TEXT,
    date TEXT NOT NULL,
    time TEXT NOT NULL,
    location TEXT NOT NULL,
    neighborhood TEXT NOT NULL
);

CREATE INDEX idx_plans_type ON plans (type);
CREATE INDEX idx_plans_neighborhood ON plans (neighborhood);
