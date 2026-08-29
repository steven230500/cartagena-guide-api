CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE businesses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,
    subtype TEXT NOT NULL,
    barrio TEXT NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL,
    image TEXT NOT NULL,
    tags TEXT[] NOT NULL DEFAULT '{}',
    description TEXT NOT NULL,
    hours TEXT,
    price_hint TEXT,
    price_typical_note TEXT,
    phone TEXT,
    web TEXT,
    email TEXT,
    instagram TEXT,
    -- Fase 5: se activa como FK a users cuando exista esa tabla.
    owner_id UUID,
    verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_businesses_type ON businesses (type);
CREATE INDEX idx_businesses_barrio ON businesses (barrio);
