CREATE TABLE user_route_progress (
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    route_id UUID NOT NULL REFERENCES routes (id) ON DELETE CASCADE,
    current_step INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, route_id)
);
