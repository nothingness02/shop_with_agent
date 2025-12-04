-- 测试数据 - 订单
INSERT INTO orders (user_id, status) VALUES
(101, 'pending'),
(102, 'completed'),
(103, 'shipped'),
(104, 'pending'),
(105, 'cancelled');

-- 测试数据 - 订单项目
INSERT INTO order_items (order_id, product_id, product_name, product_img, price, quantity, subtotal) VALUES
-- Order 1 - 用户101
(1, 1, 'iPhone 15 Pro', 'https://example.com/iphone15.jpg', 999.99, 1, 999.99),
(1, 3, 'AirPods Pro', 'https://example.com/airpods.jpg', 249.99, 2, 499.98),

-- Order 2 - 用户102
(2, 2, 'MacBook Pro 14', 'https://example.com/macbook.jpg', 1999.99, 1, 1999.99),

-- Order 3 - 用户103
(3, 5, 'Samsung Galaxy S24', 'https://example.com/s24.jpg', 899.99, 1, 899.99),
(3, 14, 'USB-C Cable 2m', 'https://example.com/cable.jpg', 19.99, 3, 59.97),

-- Order 4 - 用户104
(4, 8, 'RTX 4090 Graphics Card', 'https://example.com/gpu.jpg', 1599.99, 1, 1599.99),
(4, 9, 'Intel i9 Processor', 'https://example.com/cpu.jpg', 699.99, 1, 699.99),
(4, 10, '32GB DDR5 RAM', 'https://example.com/ram.jpg', 149.99, 2, 299.98),

-- Order 5 - 用户105
(5, 15, 'Phone Case Protective', 'https://example.com/case.jpg', 29.99, 5, 149.95);
