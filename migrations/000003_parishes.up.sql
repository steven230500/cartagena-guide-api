CREATE TABLE parishes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    neighborhood TEXT NOT NULL,
    phone TEXT
);

CREATE TABLE parish_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parish_id UUID NOT NULL REFERENCES parishes (id) ON DELETE CASCADE,
    day TEXT NOT NULL,
    times TEXT[] NOT NULL DEFAULT '{}',
    position INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_parish_schedules_parish_id ON parish_schedules (parish_id);
