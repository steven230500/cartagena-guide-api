CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_favorites (
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    business_id UUID NOT NULL REFERENCES businesses (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, business_id)
);

CREATE INDEX idx_user_favorites_user_id ON user_favorites (user_id);
