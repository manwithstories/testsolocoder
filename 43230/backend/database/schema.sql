-- 3D打印服务平台 - 数据库初始化脚本
-- PostgreSQL 14+

-- 确保 UUID 扩展可用
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 枚举类型定义
CREATE TYPE user_role AS ENUM ('designer', 'printer', 'customer', 'admin');
CREATE TYPE license_type AS ENUM ('per_purchase', 'subscription', 'commercial');
CREATE TYPE order_status AS ENUM ('pending', 'paid', 'printing', 'quality_check', 'shipped', 'delivered', 'completed', 'cancelled', 'refunded');
CREATE TYPE transaction_type AS ENUM ('income', 'expense', 'refund');
CREATE TYPE device_status AS ENUM ('idle', 'printing', 'maintenance', 'offline');
CREATE TYPE schedule_status AS ENUM ('pending', 'in_progress', 'completed', 'cancelled');
CREATE TYPE upload_status AS ENUM ('pending', 'uploading', 'completed', 'failed');

-- 用户表
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    real_name VARCHAR(50),
    avatar VARCHAR(255),
    role user_role NOT NULL,
    balance DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    credit_score DECIMAL(3,1) NOT NULL DEFAULT 5.0,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    id_card_verified BOOLEAN NOT NULL DEFAULT false,
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_email ON users(email);

-- 设计师资料表
CREATE TABLE designer_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nickname VARCHAR(50),
    bio TEXT,
    portfolio_url VARCHAR(255),
    experience_years INTEGER DEFAULT 0,
    specialties TEXT[] DEFAULT '{}',
    total_models INTEGER DEFAULT 0,
    total_sales INTEGER DEFAULT 0,
    average_rating DECIMAL(3,1) DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

-- 打印商资料表
CREATE TABLE printer_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_name VARCHAR(100),
    address TEXT,
    max_print_size VARCHAR(50),
    supported_materials TEXT[] DEFAULT '{}',
    total_orders INTEGER DEFAULT 0,
    completed_orders INTEGER DEFAULT 0,
    average_rating DECIMAL(3,1) DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

-- 材料表
CREATE TABLE materials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    type VARCHAR(30) NOT NULL,
    price_per_gram DECIMAL(8,4) NOT NULL,
    density DECIMAL(8,4) NOT NULL,
    color VARCHAR(30),
    min_layer_height DECIMAL(6,4) NOT NULL,
    max_layer_height DECIMAL(6,4) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 3D模型表
CREATE TABLE model3ds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    designer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    tags TEXT[] DEFAULT '{}',
    license_type license_type NOT NULL DEFAULT 'per_purchase',
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    subscription_price DECIMAL(10,2),
    file_format VARCHAR(10) NOT NULL,
    file_size BIGINT NOT NULL,
    original_file_path VARCHAR(255) NOT NULL,
    file_hash VARCHAR(64),
    thumbnail_url VARCHAR(255),
    volume_cm3 DECIMAL(12,4),
    bounding_box_x DECIMAL(8,2),
    bounding_box_y DECIMAL(8,2),
    bounding_box_z DECIMAL(8,2),
    estimated_print_time_min INTEGER,
    polygon_count BIGINT,
    is_printable BOOLEAN NOT NULL DEFAULT true,
    status VARCHAR(20) NOT NULL DEFAULT 'published',
    download_count INTEGER NOT NULL DEFAULT 0,
    purchase_count INTEGER NOT NULL DEFAULT 0,
    favorite_count INTEGER NOT NULL DEFAULT 0,
    average_rating DECIMAL(3,1) DEFAULT 0.0,
    total_reviews INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_model3ds_designer ON model3ds(designer_id);
CREATE INDEX idx_model3ds_category ON model3ds(category);
CREATE INDEX idx_model3ds_status ON model3ds(status);
CREATE INDEX idx_model3ds_price ON model3ds(price);

-- 模型版本表
CREATE TABLE model_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    model_id UUID NOT NULL REFERENCES model3ds(id) ON DELETE CASCADE,
    version VARCHAR(20) NOT NULL DEFAULT '1.0',
    file_path VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_hash VARCHAR(64),
    changelog TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 模型购买/订阅记录表
CREATE TABLE model_purchases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    buyer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    model_id UUID NOT NULL REFERENCES model3ds(id) ON DELETE CASCADE,
    designer_id UUID NOT NULL REFERENCES users(id),
    transaction_id UUID,
    license_type license_type NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    download_count_remaining INTEGER NOT NULL DEFAULT 5,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_model_purchases_buyer ON model_purchases(buyer_id);
CREATE INDEX idx_model_purchases_model ON model_purchases(model_id);

-- 打印订单表
CREATE TABLE print_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_no VARCHAR(32) UNIQUE NOT NULL,
    customer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    printer_id UUID REFERENCES users(id),
    model_id UUID NOT NULL REFERENCES model3ds(id) ON DELETE CASCADE,
    material_id UUID NOT NULL REFERENCES materials(id),
    material_name VARCHAR(50) NOT NULL,
    color VARCHAR(30),
    layer_height DECIMAL(6,4) NOT NULL,
    infill_percentage INTEGER NOT NULL DEFAULT 20,
    quantity INTEGER NOT NULL DEFAULT 1,
    estimated_weight_g DECIMAL(10,2) NOT NULL,
    estimated_print_time_min INTEGER NOT NULL,
    material_cost DECIMAL(10,2) NOT NULL,
    printing_cost DECIMAL(10,2) NOT NULL,
    service_fee DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    status order_status NOT NULL DEFAULT 'pending',
    quality_check_notes TEXT,
    tracking_no VARCHAR(50),
    shipping_address TEXT,
    special_instructions TEXT,
    paid_at TIMESTAMP,
    printing_started_at TIMESTAMP,
    printing_completed_at TIMESTAMP,
    quality_checked_at TIMESTAMP,
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_print_orders_customer ON print_orders(customer_id);
CREATE INDEX idx_print_orders_printer ON print_orders(printer_id);
CREATE INDEX idx_print_orders_status ON print_orders(status);
CREATE INDEX idx_print_orders_created ON print_orders(created_at);

-- 订单状态历史表
CREATE TABLE order_histories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES print_orders(id) ON DELETE CASCADE,
    status order_status NOT NULL,
    actor_id UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 分账结算表
CREATE TABLE settlements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES print_orders(id) ON DELETE CASCADE,
    total_amount DECIMAL(10,2) NOT NULL,
    platform_fee DECIMAL(10,2) NOT NULL,
    designer_share DECIMAL(10,2) NOT NULL,
    printer_share DECIMAL(10,2) NOT NULL,
    platform_fee_rate DECIMAL(5,4) NOT NULL,
    designer_share_rate DECIMAL(5,4) NOT NULL,
    printer_share_rate DECIMAL(5,4) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    settled_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 打印设备表
CREATE TABLE printer_devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    printer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    model VARCHAR(100),
    brand VARCHAR(50),
    build_volume_x DECIMAL(8,2) NOT NULL,
    build_volume_y DECIMAL(8,2) NOT NULL,
    build_volume_z DECIMAL(8,2) NOT NULL,
    nozzle_diameter DECIMAL(4,2) DEFAULT 0.4,
    supported_materials TEXT[] DEFAULT '{}',
    status device_status NOT NULL DEFAULT 'idle',
    last_maintenance_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 材料库存表
CREATE TABLE material_inventories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    printer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    material_id UUID NOT NULL REFERENCES materials(id),
    material_name VARCHAR(50) NOT NULL,
    color VARCHAR(30),
    stock_grams DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    unit_price_per_gram DECIMAL(8,4) NOT NULL,
    low_stock_threshold DECIMAL(10,2) DEFAULT 500.00,
    last_restocked_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(printer_id, material_id, color)
);

-- 排产调度表
CREATE TABLE print_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES print_orders(id) ON DELETE CASCADE,
    printer_id UUID NOT NULL REFERENCES users(id),
    device_id UUID REFERENCES printer_devices(id),
    status schedule_status NOT NULL DEFAULT 'pending',
    priority INTEGER DEFAULT 0,
    scheduled_start_at TIMESTAMP,
    scheduled_end_at TIMESTAMP,
    actual_start_at TIMESTAMP,
    actual_end_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 评价表
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES print_orders(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES users(id),
    model_id UUID NOT NULL REFERENCES model3ds(id),
    designer_id UUID NOT NULL REFERENCES users(id),
    printer_id UUID NOT NULL REFERENCES users(id),
    model_rating INTEGER NOT NULL CHECK (model_rating BETWEEN 1 AND 5),
    print_rating INTEGER NOT NULL CHECK (print_rating BETWEEN 1 AND 5),
    model_comment TEXT,
    print_comment TEXT,
    photos TEXT[] DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 收藏表
CREATE TABLE model_favorites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    model_id UUID NOT NULL REFERENCES model3ds(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, model_id)
);

-- 文件上传记录表
CREATE TABLE file_uploads (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_hash VARCHAR(64),
    file_path VARCHAR(255),
    status upload_status NOT NULL DEFAULT 'pending',
    total_chunks INTEGER,
    uploaded_chunks INTEGER DEFAULT 0,
    upload_expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 文件访问日志表
CREATE TABLE file_access_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    model_id UUID REFERENCES model3ds(id),
    file_name VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_file_access_logs_model ON file_access_logs(model_id);

-- 下载记录表
CREATE TABLE download_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    model_id UUID NOT NULL REFERENCES model3ds(id),
    purchase_id UUID REFERENCES model_purchases(id),
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT,
    ip_address VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 通知表
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    related_id UUID,
    related_type VARCHAR(50),
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_unread ON notifications(user_id, is_read);

-- 交易记录表
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_no VARCHAR(32) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type transaction_type NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    balance_after DECIMAL(12,2) NOT NULL,
    description VARCHAR(255) NOT NULL,
    related_id UUID,
    related_type VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_user ON transactions(user_id);
CREATE INDEX idx_transactions_type ON transactions(type);

-- 初始材料数据
INSERT INTO materials (name, type, price_per_gram, density, color, min_layer_height, max_layer_height, description, is_active) VALUES
('PLA 普通', 'pla', 0.05, 1.24, '白色', 0.10, 0.30, '通用打印材料，易打印，适合新手', true),
('PLA 丝绸', 'pla', 0.08, 1.24, '金色', 0.10, 0.30, '带有丝绸光泽效果的PLA材料', true),
('ABS', 'abs', 0.07, 1.04, '自然色', 0.15, 0.30, '高强度，耐冲击，需加热床和封闭打印环境', true),
('PETG', 'petg', 0.08, 1.27, '透明', 0.12, 0.30, '耐化学腐蚀，食品级，高韧性', true),
('TPU 软胶', 'tpu', 0.12, 1.20, '黑色', 0.15, 0.30, '柔性材料，适合打印垫片、轮胎等', true),
('光固化树脂 - 标准', 'resin', 0.20, 1.10, '灰色', 0.02, 0.10, '高精度光固化树脂，适合精细细节', true),
('光固化树脂 - 可浇铸', 'resin', 0.35, 1.10, '紫色', 0.02, 0.10, '可用于失蜡铸造珠宝首饰', true);

-- 更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为所有有updated_at字段的表创建触发器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_designer_profiles_updated_at ON designer_profiles;
CREATE TRIGGER update_designer_profiles_updated_at BEFORE UPDATE ON designer_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_printer_profiles_updated_at ON printer_profiles;
CREATE TRIGGER update_printer_profiles_updated_at BEFORE UPDATE ON printer_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_model3ds_updated_at ON model3ds;
CREATE TRIGGER update_model3ds_updated_at BEFORE UPDATE ON model3ds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_print_orders_updated_at ON print_orders;
CREATE TRIGGER update_print_orders_updated_at BEFORE UPDATE ON print_orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_printer_devices_updated_at ON printer_devices;
CREATE TRIGGER update_printer_devices_updated_at BEFORE UPDATE ON printer_devices
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_material_inventories_updated_at ON material_inventories;
CREATE TRIGGER update_material_inventories_updated_at BEFORE UPDATE ON material_inventories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_print_schedules_updated_at ON print_schedules;
CREATE TRIGGER update_print_schedules_updated_at BEFORE UPDATE ON print_schedules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_file_uploads_updated_at ON file_uploads;
CREATE TRIGGER update_file_uploads_updated_at BEFORE UPDATE ON file_uploads
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
