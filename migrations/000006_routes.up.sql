CREATE TABLE routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    duration TEXT NOT NULL,
    distance TEXT NOT NULL,
    difficulty TEXT NOT NULL,
    category TEXT NOT NULL,
    image TEXT NOT NULL,
    highlights TEXT[] NOT NULL DEFAULT '{}',
    audio_guide BOOLEAN NOT NULL DEFAULT false,
    offline BOOLEAN NOT NULL DEFAULT false,
    price TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE route_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES routes (id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    audio_url TEXT,
    duration TEXT,
    lat DOUBLE PRECISION,
    lng DOUBLE PRECISION,
    image TEXT,
    position INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_route_steps_route_id ON route_steps (route_id);
CREATE INDEX idx_routes_category ON routes (category);
