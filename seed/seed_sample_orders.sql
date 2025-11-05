INSERT INTO service_orders (customer_name, device_model, issue_description, status, branch_id, user_id, total_amount)
VALUES
('Budi Santoso', 'iPhone 11', 'Layar retak', 'completed',
 (SELECT id FROM branches WHERE name='Cabang Jakarta'),
 (SELECT id FROM users WHERE name='Teknisi Jakarta'),
 1200000),
('Rina Sari', 'Samsung S20', 'Tidak bisa charging', 'in_progress',
 (SELECT id FROM branches WHERE name='Cabang Bandung'),
 (SELECT id FROM users WHERE name='Teknisi Bandung'),
 850000);
