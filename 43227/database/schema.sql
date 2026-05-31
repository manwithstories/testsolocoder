-- ========================================
-- 养蜂管理与蜂蜜交易平台 数据库Schema
-- ========================================

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL CHECK (role IN ('beekeeper', 'buyer', 'inspector')),
    avatar VARCHAR(255),
    reputation DECIMAL(3,2) DEFAULT 5.00,
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS beehives (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    latitude DECIMAL(10,7) NOT NULL CHECK (latitude BETWEEN -90 AND 90),
    longitude DECIMAL(10,7) NOT NULL CHECK (longitude BETWEEN -180 AND 180),
    region VARCHAR(100),
    bee_species VARCHAR(50),
    group_name VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'harvesting')),
    health_status VARCHAR(20) DEFAULT 'healthy' CHECK (health_status IN ('healthy', 'warning', 'critical')),
    queen_status VARCHAR(20) DEFAULT 'normal' CHECK (queen_status IN ('normal', 'old', 'missing', 'new')),
    worker_count INT DEFAULT 0,
    last_inspection DATE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS health_records (
    id BIGSERIAL PRIMARY KEY,
    beehive_id BIGINT NOT NULL REFERENCES beehives(id) ON DELETE CASCADE,
    record_date DATE NOT NULL,
    queen_status VARCHAR(20) CHECK (queen_status IN ('normal', 'old', 'missing', 'new')),
    worker_count INT,
    has_disease BOOLEAN DEFAULT FALSE,
    disease_type VARCHAR(50),
    disease_severity VARCHAR(20) CHECK (disease_severity IN ('mild', 'moderate', 'severe')),
    treatment TEXT,
    season VARCHAR(20) CHECK (season IN ('spring', 'summer', 'autumn', 'winter')),
    recommendations TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS harvests (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    beehive_id BIGINT NOT NULL REFERENCES beehives(id) ON DELETE CASCADE,
    harvest_date DATE NOT NULL,
    honey_type VARCHAR(50) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity > 0),
    unit VARCHAR(10) DEFAULT 'kg',
    quality VARCHAR(20) DEFAULT 'normal' CHECK (quality IN ('normal', 'good', 'premium')),
    batch_code VARCHAR(50) UNIQUE NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS inventory (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    harvest_id BIGINT NOT NULL REFERENCES harvests(id) ON DELETE CASCADE,
    honey_type VARCHAR(50) NOT NULL,
    batch_code VARCHAR(50) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity >= 0),
    unit VARCHAR(10) DEFAULT 'kg',
    expiry_date DATE NOT NULL,
    inspection_report VARCHAR(255),
    grade VARCHAR(20) DEFAULT 'ungraded' CHECK (grade IN ('ungraded', 'grade_a', 'grade_b', 'grade_c')),
    status VARCHAR(20) DEFAULT 'in_stock' CHECK (status IN ('in_stock', 'low_stock', 'expiring_soon', 'expired', 'sold_out')),
    threshold DECIMAL(10,2) DEFAULT 10,
    price DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS inspections (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    inspector_id BIGINT REFERENCES users(id),
    inventory_id BIGINT NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    batch_code VARCHAR(50) NOT NULL,
    appointment_date DATE NOT NULL,
    inspection_date DATE,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'in_progress', 'completed', 'cancelled')),
    report_url VARCHAR(255),
    result VARCHAR(20) CHECK (result IN ('pass', 'fail', 'conditional')),
    grade VARCHAR(20) CHECK (grade IN ('grade_a', 'grade_b', 'grade_c')),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    inventory_id BIGINT NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    honey_type VARCHAR(50) NOT NULL,
    batch_code VARCHAR(50),
    price DECIMAL(10,2) NOT NULL CHECK (price > 0),
    stock DECIMAL(10,2) NOT NULL CHECK (stock >= 0),
    unit VARCHAR(10) DEFAULT 'kg',
    images TEXT[],
    grade VARCHAR(20) CHECK (grade IN ('grade_a', 'grade_b', 'grade_c')),
    status VARCHAR(20) DEFAULT 'on_sale' CHECK (status IN ('on_sale', 'off_shelf', 'sold_out')),
    view_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(50) UNIQUE NOT NULL,
    buyer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seller_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10,2) NOT NULL CHECK (unit_price > 0),
    total_amount DECIMAL(12,2) NOT NULL CHECK (total_amount > 0),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'completed', 'cancelled', 'refunded')),
    payment_status VARCHAR(20) DEFAULT 'unpaid' CHECK (payment_status IN ('unpaid', 'paid', 'refunded')),
    payment_time TIMESTAMP,
    shipping_address TEXT NOT NULL,
    tracking_number VARCHAR(100),
    tracking_status VARCHAR(50),
    buyer_rating DECIMAL(2,1),
    seller_rating DECIMAL(2,1),
    buyer_comment TEXT,
    seller_comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('disease_control', 'harvest_technique', 'seasonal_management', 'equipment', 'market', 'general')),
    tags TEXT[],
    images TEXT[],
    view_count INT DEFAULT 0,
    like_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    parent_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    like_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('disease_warning', 'seasonal_tip', 'order_status', 'inspection_result', 'low_stock', 'expiry_alert', 'system')),
    title VARCHAR(200) NOT NULL,
    content TEXT,
    related_id BIGINT,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    operation VARCHAR(50) NOT NULL,
    module VARCHAR(50),
    description TEXT,
    ip_address VARCHAR(45),
    user_agent VARCHAR(255),
    status VARCHAR(20) DEFAULT 'success' CHECK (status IN ('success', 'failed')),
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS uploads (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT,
    file_type VARCHAR(50),
    category VARCHAR(50) CHECK (category IN ('inspection_report', 'honey_image', 'avatar', 'post_image', 'other')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_beehives_user_id ON beehives(user_id);
CREATE INDEX IF NOT EXISTS idx_beehives_region ON beehives(region);
CREATE INDEX IF NOT EXISTS idx_beehives_bee_species ON beehives(bee_species);
CREATE INDEX IF NOT EXISTS idx_beehives_group_name ON beehives(group_name);
CREATE INDEX IF NOT EXISTS idx_health_records_beehive_id ON health_records(beehive_id);
CREATE INDEX IF NOT EXISTS idx_harvests_user_id ON harvests(user_id);
CREATE INDEX IF NOT EXISTS idx_harvests_beehive_id ON harvests(beehive_id);
CREATE INDEX IF NOT EXISTS idx_inventory_user_id ON inventory(user_id);
CREATE INDEX IF NOT EXISTS idx_inventory_batch_code ON inventory(batch_code);
CREATE INDEX IF NOT EXISTS idx_inventory_status ON inventory(status);
CREATE INDEX IF NOT EXISTS idx_products_user_id ON products(user_id);
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
CREATE INDEX IF NOT EXISTS idx_orders_buyer_id ON orders(buyer_id);
CREATE INDEX IF NOT EXISTS idx_orders_seller_id ON orders(seller_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_category ON posts(category);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_operation_logs_user_id ON operation_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_created_at ON operation_logs(created_at);
