-- Pet Boarding & Daycare Management Platform - Database Schema
-- PostgreSQL Migration

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20),
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'owner',
    avatar_url VARCHAR(500),
    real_name VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS store_infos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    store_name VARCHAR(100) NOT NULL,
    address VARCHAR(500),
    license_no VARCHAR(50),
    business_hours VARCHAR(200),
    description VARCHAR(1000),
    rating FLOAT DEFAULT 0,
    review_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS keeper_infos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    store_id UUID REFERENCES users(id),
    nick_name VARCHAR(50),
    experience INT DEFAULT 0,
    specialty VARCHAR(500),
    rating FLOAT DEFAULT 0,
    review_count INT DEFAULT 0,
    certification VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(50) NOT NULL,
    species VARCHAR(20) NOT NULL,
    breed VARCHAR(50),
    gender VARCHAR(10),
    birth_date DATE,
    weight FLOAT DEFAULT 0,
    color VARCHAR(30),
    avatar_url VARCHAR(500),
    allergies VARCHAR(500),
    diet_habit VARCHAR(500),
    temperament VARCHAR(200),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS vaccine_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id),
    vaccine_name VARCHAR(100) NOT NULL,
    vaccinated_at TIMESTAMP NOT NULL,
    expire_at TIMESTAMP NOT NULL,
    hospital VARCHAR(100),
    proof_url VARCHAR(500),
    is_valid BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS deworm_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id),
    deworm_type VARCHAR(50) NOT NULL,
    dewormed_at TIMESTAMP NOT NULL,
    expire_at TIMESTAMP NOT NULL,
    medicine VARCHAR(100),
    is_valid BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS boarding_packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    description VARCHAR(1000),
    price_per_day FLOAT NOT NULL,
    capacity INT DEFAULT 1,
    features VARCHAR(1000),
    is_available BOOLEAN DEFAULT true,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_no VARCHAR(50) UNIQUE NOT NULL,
    owner_id UUID NOT NULL REFERENCES users(id),
    pet_id UUID NOT NULL REFERENCES pets(id),
    store_id UUID NOT NULL REFERENCES users(id),
    package_id UUID NOT NULL REFERENCES boarding_packages(id),
    package_type VARCHAR(20) NOT NULL,
    check_in_date TIMESTAMP NOT NULL,
    check_out_date TIMESTAMP NOT NULL,
    total_days INT NOT NULL,
    total_amount FLOAT NOT NULL,
    deposit_amount FLOAT DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    keeper_id UUID REFERENCES users(id),
    remark VARCHAR(500),
    cancel_reason VARCHAR(500),
    cancelled_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reservations_store_date ON reservations(store_id, check_in_date, check_out_date);
CREATE INDEX IF NOT EXISTS idx_reservations_status ON reservations(status);

CREATE TABLE IF NOT EXISTS daily_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reservation_id UUID NOT NULL REFERENCES reservations(id),
    pet_id UUID NOT NULL REFERENCES pets(id),
    keeper_id UUID NOT NULL REFERENCES users(id),
    record_date TIMESTAMP NOT NULL,
    feed_status VARCHAR(200),
    activity VARCHAR(500),
    health_status VARCHAR(500),
    mood VARCHAR(200),
    photos VARCHAR(2000),
    remark VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_daily_records_reservation ON daily_records(reservation_id);
CREATE INDEX IF NOT EXISTS idx_daily_records_pet ON daily_records(pet_id);

CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reservation_id UUID NOT NULL UNIQUE REFERENCES reservations(id),
    owner_id UUID NOT NULL REFERENCES users(id),
    store_id UUID NOT NULL REFERENCES users(id),
    keeper_id UUID REFERENCES users(id),
    store_rating INT NOT NULL CHECK (store_rating >= 1 AND store_rating <= 5),
    keeper_rating INT CHECK (keeper_rating >= 1 AND keeper_rating <= 5),
    content VARCHAR(1000),
    images VARCHAR(2000),
    reply VARCHAR(1000),
    reply_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reviews_store ON reviews(store_id);
CREATE INDEX IF NOT EXISTS idx_reviews_keeper ON reviews(keeper_id);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_no VARCHAR(50) UNIQUE NOT NULL,
    reservation_id UUID NOT NULL REFERENCES reservations(id),
    owner_id UUID NOT NULL REFERENCES users(id),
    store_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(20) NOT NULL,
    amount FLOAT NOT NULL,
    pay_status VARCHAR(20) NOT NULL DEFAULT 'unpaid',
    pay_method VARCHAR(30),
    transaction_id VARCHAR(100),
    paid_at TIMESTAMP,
    refund_amount FLOAT DEFAULT 0,
    refund_at TIMESTAMP,
    remark VARCHAR(500),
    amount_hash VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orders_reservation ON orders(reservation_id);
CREATE INDEX IF NOT EXISTS idx_orders_store ON orders(store_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(pay_status);

CREATE TABLE IF NOT EXISTS health_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    pet_id UUID NOT NULL REFERENCES pets(id),
    alert_type VARCHAR(30) NOT NULL,
    title VARCHAR(200) NOT NULL,
    content VARCHAR(500),
    record_id UUID,
    expire_at TIMESTAMP,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_alerts_user ON health_alerts(user_id);
CREATE INDEX IF NOT EXISTS idx_alerts_read ON health_alerts(user_id, is_read);

CREATE TABLE IF NOT EXISTS operation_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    username VARCHAR(50),
    action VARCHAR(100) NOT NULL,
    method VARCHAR(10),
    url VARCHAR(500),
    ip VARCHAR(50),
    params TEXT,
    result TEXT,
    status INT,
    exec_time BIGINT,
    error_msg VARCHAR(1000),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_logs_user ON operation_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_logs_created ON operation_logs(created_at);
