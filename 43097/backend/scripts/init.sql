-- 酒店管理系统数据库初始化脚本
-- PostgreSQL 13+

-- 创建数据库（如果不存在需要手动执行）
-- CREATE DATABASE hotel_db;

-- 连接到数据库
\c hotel_db;

-- 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ========== 用户表 ==========
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    real_name VARCHAR(100),
    phone VARCHAR(20),
    email VARCHAR(100),
    role VARCHAR(20) NOT NULL DEFAULT 'user', -- admin, frontdesk, user
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, inactive
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);

-- ========== 会员等级表 ==========
CREATE TABLE IF NOT EXISTS member_levels (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    discount_rate DECIMAL(5,2) NOT NULL DEFAULT 1.00,
    points_rate DECIMAL(5,2) NOT NULL DEFAULT 1.00,
    min_points INT NOT NULL DEFAULT 0,
    max_points INT NOT NULL DEFAULT 999999,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ========== 会员表 ==========
CREATE TABLE IF NOT EXISTS members (
    id BIGSERIAL PRIMARY KEY,
    member_no VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(100),
    id_card VARCHAR(20),
    level_id BIGINT REFERENCES member_levels(id),
    points INT NOT NULL DEFAULT 0,
    balance DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, inactive
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_members_member_no ON members(member_no);
CREATE INDEX idx_members_phone ON members(phone);
CREATE INDEX idx_members_level_id ON members(level_id);

-- ========== 会员积分流水表 ==========
CREATE TABLE IF NOT EXISTS member_points_logs (
    id BIGSERIAL PRIMARY KEY,
    member_id BIGINT NOT NULL REFERENCES members(id),
    type VARCHAR(20) NOT NULL, -- earn, use, recharge, refund
    points INT NOT NULL,
    balance_before INT NOT NULL,
    balance_after INT NOT NULL,
    description TEXT,
    related_type VARCHAR(50),
    related_id BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_points_logs_member_id ON member_points_logs(member_id);
CREATE INDEX idx_points_logs_type ON member_points_logs(type);

-- ========== 房型表 ==========
CREATE TABLE IF NOT EXISTS room_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    base_price DECIMAL(10,2) NOT NULL,
    bed_count INT NOT NULL DEFAULT 1,
    max_guests INT NOT NULL DEFAULT 2,
    area INT,
    facilities JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ========== 房间表 ==========
CREATE TABLE IF NOT EXISTS rooms (
    id BIGSERIAL PRIMARY KEY,
    room_no VARCHAR(20) UNIQUE NOT NULL,
    floor INT NOT NULL,
    room_type_id BIGINT NOT NULL REFERENCES room_types(id),
    status VARCHAR(20) NOT NULL DEFAULT 'available', -- available, occupied, reserved, maintenance
    price DECIMAL(10,2) NOT NULL,
    facilities JSONB,
    remark TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_rooms_room_no ON rooms(room_no);
CREATE INDEX idx_rooms_floor ON rooms(floor);
CREATE INDEX idx_rooms_status ON rooms(status);
CREATE INDEX idx_rooms_type ON rooms(room_type_id);

-- ========== 预订表 ==========
CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    booking_no VARCHAR(50) UNIQUE NOT NULL,
    room_id BIGINT NOT NULL REFERENCES rooms(id),
    member_id BIGINT REFERENCES members(id),
    guest_name VARCHAR(100) NOT NULL,
    guest_phone VARCHAR(20) NOT NULL,
    guest_id_card VARCHAR(20),
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,
    days INT NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, confirmed, cancelled, completed, no_show
    paid_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    remarks TEXT,
    cancel_deadline TIMESTAMP,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_bookings_booking_no ON bookings(booking_no);
CREATE INDEX idx_bookings_room_id ON bookings(room_id);
CREATE INDEX idx_bookings_member_id ON bookings(member_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_dates ON bookings(check_in_date, check_out_date);

-- ========== 入住表 ==========
CREATE TABLE IF NOT EXISTS check_ins (
    id BIGSERIAL PRIMARY KEY,
    check_in_no VARCHAR(50) UNIQUE NOT NULL,
    booking_id BIGINT REFERENCES bookings(id),
    room_id BIGINT NOT NULL REFERENCES rooms(id),
    guest_name VARCHAR(100) NOT NULL,
    guest_phone VARCHAR(20) NOT NULL,
    guest_id_card VARCHAR(20),
    check_in_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expected_check_out DATE NOT NULL,
    actual_check_out TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, checked_out
    deposit DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    extra_charges DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    remarks TEXT,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_checkins_checkin_no ON check_ins(check_in_no);
CREATE INDEX idx_checkins_room_id ON check_ins(room_id);
CREATE INDEX idx_checkins_booking_id ON check_ins(booking_id);
CREATE INDEX idx_checkins_status ON check_ins(status);
CREATE INDEX idx_checkins_checkin_time ON check_ins(check_in_time);

-- ========== 支付记录表 ==========
CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    payment_no VARCHAR(50) UNIQUE NOT NULL,
    order_type VARCHAR(20) NOT NULL, -- booking, checkin
    order_id BIGINT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    payment_method VARCHAR(20) NOT NULL, -- cash, wechat, alipay, card, transfer
    payment_type VARCHAR(20) NOT NULL, -- prepaid, extra, refund, deposit
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, completed, failed, refunded
    transaction_id VARCHAR(100),
    remark TEXT,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payments_payment_no ON payments(payment_no);
CREATE INDEX idx_payments_order ON payments(order_type, order_id);
CREATE INDEX idx_payments_method ON payments(payment_method);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at);

-- ========== 初始化数据 ==========

-- 初始化会员等级
INSERT INTO member_levels (name, discount_rate, points_rate, min_points, max_points, description) VALUES
('普通会员', 1.00, 1.00, 0, 999, '注册即成为普通会员'),
('银卡会员', 0.95, 1.20, 1000, 4999, '享95折优惠，积分1.2倍'),
('金卡会员', 0.90, 1.50, 5000, 19999, '享9折优惠，积分1.5倍'),
('钻石会员', 0.85, 2.00, 20000, 999999, '享85折优惠，积分2倍')
ON CONFLICT DO NOTHING;

-- 初始化管理员用户 (密码: admin123, 需要使用bcrypt哈希)
-- 实际密码哈希需要在应用启动时通过seed生成

-- 初始化房型
INSERT INTO room_types (name, description, base_price, bed_count, max_guests, area, facilities) VALUES
('标准单人间', '经济实惠的单人间，适合单人入住', 198.00, 1, 1, 20, '["wifi","tv","air-conditioner","24h-hot-water"]'),
('标准双床间', '双床标准间，适合朋友或同事入住', 258.00, 2, 2, 25, '["wifi","tv","air-conditioner","24h-hot-water"]'),
('豪华大床房', '宽敞舒适的大床房，配备高品质设施', 388.00, 1, 2, 35, '["wifi","tv","air-conditioner","24h-hot-water","minibar","safe"]'),
('行政套房', '豪华套房，独立客厅和卧室，尊享行政礼遇', 688.00, 1, 2, 60, '["wifi","tv","air-conditioner","24h-hot-water","minibar","safe","bathtub","living-room"]')
ON CONFLICT DO NOTHING;

-- 初始化房间数据
INSERT INTO rooms (room_no, floor, room_type_id, status, price, remark) VALUES
('101', 1, 1, 'available', 198.00, ''),
('102', 1, 1, 'available', 198.00, ''),
('103', 1, 2, 'available', 258.00, ''),
('104', 1, 2, 'available', 258.00, ''),
('201', 2, 1, 'available', 198.00, ''),
('202', 2, 1, 'available', 198.00, ''),
('203', 2, 2, 'available', 258.00, ''),
('204', 2, 2, 'available', 258.00, ''),
('301', 3, 3, 'available', 388.00, ''),
('302', 3, 3, 'available', 388.00, ''),
('401', 4, 4, 'available', 688.00, ''),
('402', 4, 4, 'maintenance', 688.00, '维护中')
ON CONFLICT DO NOTHING;
