ALTER TABLE businesses
    ADD CONSTRAINT businesses_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE SET NULL;

CREATE TABLE business_claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL REFERENCES businesses (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    message TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_business_claims_pending_unique
    ON business_claims (business_id, user_id) WHERE status = 'pending';
CREATE INDEX idx_business_claims_status ON business_claims (status);
CREATE INDEX idx_business_claims_user_id ON business_claims (user_id);
