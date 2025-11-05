-- Indexes untuk performa query
CREATE INDEX idx_users_branch_id ON users(branch_id);
CREATE INDEX idx_orders_branch_id ON service_orders(branch_id);
CREATE INDEX idx_orders_user_id ON service_orders(user_id);
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_created_at ON payments(created_at);
