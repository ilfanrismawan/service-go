CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES service_orders(id) ON DELETE CASCADE,
    amount NUMERIC(12,2) NOT NULL,
    method VARCHAR(30),
    status VARCHAR(30) DEFAULT 'completed',
    created_at TIMESTAMP DEFAULT NOW()
);
