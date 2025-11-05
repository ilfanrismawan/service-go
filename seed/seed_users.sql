INSERT INTO users (name, email, password_hash, role, branch_id)
VALUES
('Admin Utama', 'admin@service.com', '$2a$10$xxxx', 'ADMIN', NULL),
('Teknisi Jakarta', 'teknisi.jkt@service.com', '$2a$10$xxxx', 'TECHNICIAN', (SELECT id FROM branches WHERE name='Cabang Jakarta')),
('Teknisi Bandung', 'teknisi.bdg@service.com', '$2a$10$xxxx', 'TECHNICIAN', (SELECT id FROM branches WHERE name='Cabang Bandung'));
