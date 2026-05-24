-- 企业工商注册代办平台数据库初始化脚本

-- 创建数据库
CREATE DATABASE IF NOT EXISTS business_registration DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE business_registration;

-- 初始化费用标准
INSERT INTO fee_standards (company_type, naming_fee, registration_fee, tax_fee, bank_fee, seal_fee, service_fee, capital_rate, created_at, updated_at) VALUES
('llc', 100, 500, 200, 300, 150, 500, 0.5, NOW(), NOW()),
('joint_stock', 200, 1000, 300, 500, 200, 1000, 1.0, NOW(), NOW()),
('sole', 50, 200, 100, 150, 80, 300, 0.3, NOW(), NOW()),
('partnership', 80, 300, 150, 200, 100, 400, 0.4, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 初始化通知模板
INSERT INTO notification_templates (code, name, type, title, content, variables, is_active, created_at, updated_at) VALUES
('APPLICATION_STATUS_CHANGE', '申请状态变更通知', 'system', '{{company_name}}的申请状态已变更', '尊敬的用户，您的公司{{company_name}}（申请编号：{{application_no}}）状态已变更为：{{status}}', 'company_name,application_no,status', 1, NOW(), NOW()),
('STEP_COMPLETED', '环节完成通知', 'system', '{{company_name}}的{{step_name}}已完成', '尊敬的用户，您的公司{{company_name}}（申请编号：{{application_no}}）的{{step_name}}环节已完成', 'company_name,application_no,step_name', 1, NOW(), NOW()),
('NEW_APPLICATION', '新申请通知', 'system', '有新的申请待处理', '您有一个新的公司注册申请待处理，请及时查看', '', 1, NOW(), NOW()),
('PAYMENT_SUCCESS', '支付成功通知', 'system', '支付成功', '您的申请费用已支付成功，申请即将进入审核流程', '', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 初始化优惠策略
INSERT INTO discount_policies (name, code, type, value, min_amount, max_discount, start_date, end_date, is_active, created_at, updated_at) VALUES
('新用户优惠', 'NEWUSER100', 'fixed', 100, 1000, 100, NULL, NULL, 1, NOW(), NOW()),
('春季促销', 'SPRING20', 'percent', 20, 2000, 500, NOW(), DATE_ADD(NOW(), INTERVAL 3 MONTH), 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 初始化管理员账户（密码: admin123）
INSERT INTO users (username, password, real_name, email, phone, role, status, created_at, updated_at) VALUES
('admin', '$2a$14$wVsaPvJnJJsomWArouWCtusem6S/.Gauq/GjOEI9H8VW2nRO0bKbK', '系统管理员', 'admin@example.com', '13800000000', 'admin', 'active', NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();
