-- 查询测试数据的 SQL 脚本
-- 可用于验证测试数据是否正确导入

-- 查看所有商店
-- SELECT * FROM shops WHERE deleted_at IS NULL ORDER BY id;

-- 查看特定商店的所有商品
-- SELECT p.* FROM products p
-- WHERE p.shop_id = 1 AND p.deleted_at IS NULL
-- ORDER BY p.id;

-- 查看商品分类信息
-- SELECT p.id, p.name, p.price, c.name as category
-- FROM products p
-- LEFT JOIN product_categories pc ON p.id = pc.product_id
-- LEFT JOIN categories c ON pc.category_id = c.id
-- WHERE p.deleted_at IS NULL
-- ORDER BY p.id;

-- 查看所有订单及其项目
-- SELECT o.id as order_id, o.user_id, o.status, 
--        COUNT(oi.id) as item_count, SUM(oi.subtotal) as total_amount
-- FROM orders o
-- LEFT JOIN order_items oi ON o.id = oi.order_id
-- WHERE o.deleted_at IS NULL
-- GROUP BY o.id, o.user_id, o.status
-- ORDER BY o.id;

-- 查看特定订单的详细信息
-- SELECT o.id, o.user_id, o.status, o.created_at,
--        oi.product_name, oi.quantity, oi.price, oi.subtotal
-- FROM orders o
-- LEFT JOIN order_items oi ON o.id = oi.order_id
-- WHERE o.id = 1 AND o.deleted_at IS NULL;

-- 统计各个商店的商品数量和库存
-- SELECT s.id, s.name, COUNT(p.id) as product_count, 
--        SUM(p.stock) as total_stock,
--        AVG(p.price) as avg_price
-- FROM shops s
-- LEFT JOIN products p ON s.id = p.shop_id AND p.deleted_at IS NULL
-- WHERE s.deleted_at IS NULL
-- GROUP BY s.id, s.name
-- ORDER BY s.id;

-- 查找特定名称的商品
-- SELECT * FROM products 
-- WHERE name LIKE '%iPhone%' AND deleted_at IS NULL;

-- 查找价格在某个范围内的商品
-- SELECT * FROM products
-- WHERE price BETWEEN 100 AND 500 AND deleted_at IS NULL
-- ORDER BY price DESC;

-- 查找库存较低的商品（库存少于20件）
-- SELECT s.name as shop_name, p.name, p.stock, p.price
-- FROM products p
-- JOIN shops s ON p.shop_id = s.id
-- WHERE p.stock < 20 AND p.deleted_at IS NULL AND s.deleted_at IS NULL
-- ORDER BY p.stock ASC;

-- 查看每个用户的订单总额
-- SELECT o.user_id, COUNT(o.id) as order_count, 
--        SUM(oi.subtotal) as total_spent
-- FROM orders o
-- LEFT JOIN order_items oi ON o.id = oi.order_id
-- WHERE o.deleted_at IS NULL
-- GROUP BY o.user_id
-- ORDER BY total_spent DESC;
