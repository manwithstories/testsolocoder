-- 添加用户角色字段并创建默认管理员账号
-- 执行此脚本前请先备份数据库

BEGIN;

-- 添加角色字段（如果不存在）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'role'
    ) THEN
        ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user';
    END IF;
END $$;

-- 创建默认管理员账号（如果不存在）
-- 邮箱: admin@example.com
-- 密码: admin123456
-- 请在生产环境中及时修改默认密码！
INSERT INTO users (username, email, password, role, credit_score, is_banned, created_at, updated_at)
SELECT 'admin', 'admin@example.com', 
       '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
       'admin', 100, false, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin@example.com');

COMMIT;
