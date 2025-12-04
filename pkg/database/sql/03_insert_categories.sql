-- 测试数据 - 分类
INSERT INTO categories (name, description) VALUES
('Smartphones', 'Mobile phones and devices'),
('Laptops', 'Portable computers'),
('Accessories', 'Phone and computer accessories'),
('Components', 'Computer hardware components'),
('Home Appliances', 'Smart home devices'),
('Gaming', 'Gaming peripherals and equipment');

-- 关联产品和分类
INSERT INTO product_categories (product_id, category_id) VALUES
-- iPhone products
(1, 1),  -- iPhone 15 Pro -> Smartphones
(3, 3),  -- AirPods Pro -> Accessories
(4, 1),  -- Apple Watch Series 9 -> Smartphones

-- Samsung products
(5, 1),  -- Samsung Galaxy S24 -> Smartphones
(6, 5),  -- Samsung 4K TV -> Home Appliances
(7, 5),  -- Samsung Washing Machine -> Home Appliances

-- Laptops
(2, 2),  -- MacBook Pro 14 -> Laptops

-- Components
(8, 4),  -- RTX 4090 -> Components
(9, 4),  -- Intel i9 -> Components
(10, 4), -- 32GB DDR5 RAM -> Components

-- Gaming/Peripherals
(11, 6), -- Mechanical Keyboard RGB -> Gaming
(12, 6), -- Wireless Mouse Pro -> Gaming
(13, 3), -- 4K Monitor -> Accessories

-- Accessories
(14, 3), -- USB-C Cable -> Accessories
(15, 3), -- Phone Case -> Accessories
(16, 3); -- Screen Protector -> Accessories
