-- 插入测试数据 - 商店
INSERT INTO shops (name, description, owner_id) VALUES
('Apple Store', 'Official Apple retail store', 1),
('Samsung Electronics', 'Samsung authorized dealer', 2),
('Tech World', 'Electronics and gadgets', 3),
('Digital Hub', 'Computers and peripherals', 4),
('Mobile Paradise', 'Smartphones and accessories', 5);

-- 插入测试数据 - 商品
INSERT INTO products (shop_id, name, description, price, stock, product_img) VALUES
-- Apple Store products
(1, 'iPhone 15 Pro', 'Latest Apple flagship smartphone', 999.99, 50, 'https://example.com/iphone15.jpg'),
(1, 'MacBook Pro 14', 'High-performance laptop', 1999.99, 20, 'https://example.com/macbook.jpg'),
(1, 'AirPods Pro', 'Premium wireless earbuds', 249.99, 100, 'https://example.com/airpods.jpg'),
(1, 'Apple Watch Series 9', 'Smart wearable', 399.99, 30, 'https://example.com/watch.jpg'),

-- Samsung Electronics products
(2, 'Samsung Galaxy S24', 'Latest Samsung flagship', 899.99, 40, 'https://example.com/s24.jpg'),
(2, 'Samsung 4K TV 55"', 'Ultra HD television', 799.99, 15, 'https://example.com/tv55.jpg'),
(2, 'Samsung Washing Machine', 'Smart washing machine', 599.99, 10, 'https://example.com/washer.jpg'),

-- Tech World products
(3, 'RTX 4090 Graphics Card', 'High-end gaming GPU', 1599.99, 5, 'https://example.com/gpu.jpg'),
(3, 'Intel i9 Processor', 'High-performance CPU', 699.99, 12, 'https://example.com/cpu.jpg'),
(3, '32GB DDR5 RAM', 'Fast memory modules', 149.99, 25, 'https://example.com/ram.jpg'),

-- Digital Hub products
(4, 'Mechanical Keyboard RGB', 'Gaming keyboard with RGB lights', 149.99, 45, 'https://example.com/keyboard.jpg'),
(4, 'Wireless Mouse Pro', 'Ergonomic wireless mouse', 79.99, 60, 'https://example.com/mouse.jpg'),
(4, '4K Monitor 32"', 'Ultra-wide display', 499.99, 8, 'https://example.com/monitor.jpg'),

-- Mobile Paradise products
(5, 'USB-C Cable 2m', 'High-speed charging cable', 19.99, 200, 'https://example.com/cable.jpg'),
(5, 'Phone Case Protective', 'Durable protective case', 29.99, 150, 'https://example.com/case.jpg'),
(5, 'Screen Protector Pack', 'Tempered glass screen protector', 14.99, 300, 'https://example.com/protector.jpg');
