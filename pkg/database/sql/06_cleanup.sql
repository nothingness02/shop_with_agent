-- 清空测试数据脚本
-- 用于重置数据库，按照外键依赖顺序删除

-- 删除订单相关数据
DELETE FROM order_items WHERE order_id IN (SELECT id FROM orders);
DELETE FROM orders;

-- 删除商品相关数据
DELETE FROM product_categories;
DELETE FROM products;
DELETE FROM categories;

-- 删除商店数据
DELETE FROM shops;

-- 重置自增ID序列（PostgreSQL 语法）
ALTER SEQUENCE orders_id_seq RESTART WITH 1;
ALTER SEQUENCE order_items_id_seq RESTART WITH 1;
ALTER SEQUENCE products_id_seq RESTART WITH 1;
ALTER SEQUENCE categories_id_seq RESTART WITH 1;
ALTER SEQUENCE shops_id_seq RESTART WITH 1;
