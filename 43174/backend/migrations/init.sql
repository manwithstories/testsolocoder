-- 创建数据库
CREATE DATABASE campus_trade;

-- 连接到数据库
\c campus_trade;

-- 启用UUID扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL DEFAULT 'student',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    real_name VARCHAR(50),
    school_name VARCHAR(100),
    student_id VARCHAR(50),
    student_card_url VARCHAR(500),
    business_license VARCHAR(500),
    avatar VARCHAR(500),
    rating DOUBLE PRECISION DEFAULT 5.0,
    rating_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建分类表
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    parent_id UUID,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES categories(id)
);

-- 创建教材表
CREATE TABLE IF NOT EXISTS textbooks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    isbn VARCHAR(20),
    title VARCHAR(200) NOT NULL,
    author VARCHAR(100),
    course_name VARCHAR(200),
    edition VARCHAR(50),
    publisher VARCHAR(100),
    original_price DOUBLE PRECISION,
    price DOUBLE PRECISION NOT NULL,
    condition VARCHAR(20),
    description TEXT,
    cover_image VARCHAR(500),
    status VARCHAR(20) DEFAULT 'available',
    seller_id UUID REFERENCES users(id),
    category_id UUID REFERENCES categories(id),
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建笔记表
CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(200) NOT NULL,
    subject VARCHAR(100),
    course_name VARCHAR(200),
    description TEXT,
    file_url VARCHAR(500) NOT NULL,
    file_type VARCHAR(20),
    file_size BIGINT,
    cover_image VARCHAR(500),
    uploader_id UUID REFERENCES users(id),
    category_id UUID REFERENCES categories(id),
    download_count INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    rating DOUBLE PRECISION DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建交易表
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    textbook_id UUID REFERENCES textbooks(id),
    seller_id UUID REFERENCES users(id),
    buyer_id UUID REFERENCES users(id),
    type VARCHAR(20),
    agreed_price DOUBLE PRECISION,
    status VARCHAR(20) DEFAULT 'pending',
    exchange_item VARCHAR(500),
    negotiation_history TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建订单表
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_no VARCHAR(50) UNIQUE NOT NULL,
    buyer_id UUID REFERENCES users(id),
    seller_id UUID REFERENCES users(id),
    total_amount DOUBLE PRECISION NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    payment_method VARCHAR(50),
    payment_status VARCHAR(20),
    transaction_id UUID REFERENCES transactions(id),
    shipping_address VARCHAR(500),
    tracking_number VARCHAR(50),
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建订单项目表
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID REFERENCES orders(id),
    textbook_id UUID REFERENCES textbooks(id),
    quantity INTEGER DEFAULT 1,
    price DOUBLE PRECISION,
    subtotal DOUBLE PRECISION,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建订单状态历史表
CREATE TABLE IF NOT EXISTS order_status_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID REFERENCES orders(id),
    status VARCHAR(20),
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建消息表
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id UUID REFERENCES users(id),
    receiver_id UUID REFERENCES users(id),
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    related_order_id UUID REFERENCES orders(id),
    is_dispute BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建评价表
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    target_type VARCHAR(20),
    textbook_id UUID REFERENCES textbooks(id),
    note_id UUID REFERENCES notes(id),
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    content TEXT,
    is_hidden BOOLEAN DEFAULT FALSE,
    is_malicious BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建通知表
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    type VARCHAR(50),
    title VARCHAR(200),
    content TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_textbooks_isbn ON textbooks(isbn);
CREATE INDEX IF NOT EXISTS idx_textbooks_course_name ON textbooks(course_name);
CREATE INDEX IF NOT EXISTS idx_textbooks_seller_id ON textbooks(seller_id);
CREATE INDEX IF NOT EXISTS idx_textbooks_status ON textbooks(status);
CREATE INDEX IF NOT EXISTS idx_notes_subject ON notes(subject);
CREATE INDEX IF NOT EXISTS idx_notes_uploader_id ON notes(uploader_id);
CREATE INDEX IF NOT EXISTS idx_orders_buyer_id ON orders(buyer_id);
CREATE INDEX IF NOT EXISTS idx_orders_seller_id ON orders(seller_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages(receiver_id);
CREATE INDEX IF NOT EXISTS idx_reviews_textbook_id ON reviews(textbook_id);
CREATE INDEX IF NOT EXISTS idx_reviews_note_id ON reviews(note_id);

-- 插入初始数据
INSERT INTO categories (name, sort_order) VALUES
    ('计算机科学', 1),
    ('数学', 2),
    ('物理', 3),
    ('化学', 4),
    ('生物', 5),
    ('经济', 6),
    ('管理', 7),
    ('文学', 8),
    ('历史', 9),
    ('哲学', 10);

-- 创建管理员账户 (密码: admin123)
INSERT INTO users (username, email, password, role, status, real_name)
VALUES (
    'admin',
    'admin@campus.edu',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'admin',
    'active',
    '系统管理员'
);
