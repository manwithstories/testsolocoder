-- ============================================================
-- Matchmaking Platform Database Schema
-- ============================================================

CREATE DATABASE IF NOT EXISTS matchmaking DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE matchmaking;

-- -----------------------------------------------------------
-- Users Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username        VARCHAR(50)  NOT NULL UNIQUE,
    password        VARCHAR(255) NOT NULL,
    phone           VARCHAR(20)  UNIQUE,
    email           VARCHAR(100) UNIQUE,
    role            VARCHAR(20)  NOT NULL DEFAULT 'user',
    status          VARCHAR(20)  NOT NULL DEFAULT 'active',
    verify_status   VARCHAR(20)  NOT NULL DEFAULT 'pending',
    real_name       VARCHAR(50),
    id_card         VARCHAR(20),
    id_card_front   VARCHAR(255),
    id_card_back    VARCHAR(255),
    avatar          VARCHAR(255),
    member_level    VARCHAR(20)  NOT NULL DEFAULT 'free',
    member_expire   DATETIME     NULL,
    last_login_at   DATETIME     NULL,
    last_login_ip   VARCHAR(50),
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      DATETIME     NULL,
    INDEX idx_role (role),
    INDEX idx_status (status),
    INDEX idx_verify (verify_status),
    INDEX idx_member_level (member_level)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- User Profiles Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS profiles (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL UNIQUE,
    nickname        VARCHAR(50),
    gender          VARCHAR(10),
    birthday        DATE,
    age             INT          DEFAULT 0,
    height          INT          DEFAULT 0,
    weight          INT          DEFAULT 0,
    education       VARCHAR(20),
    occupation      VARCHAR(100),
    income          VARCHAR(20),
    city            VARCHAR(50),
    district        VARCHAR(50),
    address         VARCHAR(255),
    latitude        DOUBLE       DEFAULT 0,
    longitude       DOUBLE       DEFAULT 0,
    intro           TEXT,
    hobbies         TEXT,
    tags            TEXT,
    photos          TEXT,
    min_age         INT          DEFAULT 0,
    max_age         INT          DEFAULT 0,
    min_height      INT          DEFAULT 0,
    max_height      INT          DEFAULT 0,
    prefer_education VARCHAR(20),
    prefer_income   VARCHAR(20),
    prefer_city     VARCHAR(50),
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_gender (gender),
    INDEX idx_city (city),
    INDEX idx_age (age)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Match Records Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS match_records (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL,
    target_id       BIGINT UNSIGNED NOT NULL,
    match_score     DOUBLE       DEFAULT 0,
    match_reason    TEXT,
    is_favorited    TINYINT(1)   DEFAULT 0,
    is_blocked      TINYINT(1)   DEFAULT 0,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_target_id (target_id),
    INDEX idx_user_target (user_id, target_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Date Records Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS date_records (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    initiator_id    BIGINT UNSIGNED NOT NULL,
    receiver_id     BIGINT UNSIGNED NOT NULL,
    matchmaker_id   BIGINT UNSIGNED NULL,
    title           VARCHAR(200) NOT NULL,
    location        VARCHAR(200),
    date_at         DATETIME     NOT NULL,
    duration        INT          DEFAULT 60,
    status          VARCHAR(20)  NOT NULL DEFAULT 'pending',
    note            TEXT,
    reminded        TINYINT(1)   DEFAULT 0,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_initiator (initiator_id),
    INDEX idx_receiver (receiver_id),
    INDEX idx_matchmaker (matchmaker_id),
    INDEX idx_status (status),
    INDEX idx_date_at (date_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Date Reviews Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS date_reviews (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    date_id         BIGINT UNSIGNED NOT NULL,
    reviewer_id     BIGINT UNSIGNED NOT NULL,
    target_id       BIGINT UNSIGNED NOT NULL,
    rating          INT          NOT NULL,
    content         TEXT,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_date_id (date_id),
    INDEX idx_reviewer (reviewer_id),
    INDEX idx_target (target_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Matchmaker Members Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS matchmaker_members (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    matchmaker_id   BIGINT UNSIGNED NOT NULL,
    member_id       BIGINT UNSIGNED NOT NULL,
    status          VARCHAR(20)  NOT NULL DEFAULT 'active',
    joined_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_matchmaker_member (matchmaker_id, member_id),
    INDEX idx_matchmaker (matchmaker_id),
    INDEX idx_member (member_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Matchmaker Services Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS matchmaker_services (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    matchmaker_id   BIGINT UNSIGNED NOT NULL,
    member_a_id     BIGINT UNSIGNED NOT NULL,
    member_b_id     BIGINT UNSIGNED NOT NULL,
    date_id         BIGINT UNSIGNED NULL,
    service_type    VARCHAR(50)  NOT NULL,
    note            TEXT,
    status          VARCHAR(20)  NOT NULL DEFAULT 'progress',
    progress        INT          DEFAULT 0,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_matchmaker (matchmaker_id),
    INDEX idx_member_a (member_a_id),
    INDEX idx_member_b (member_b_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Matchmaker Stats Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS matchmaker_stats (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    matchmaker_id   BIGINT UNSIGNED NOT NULL UNIQUE,
    total_members   INT          DEFAULT 0,
    total_services  INT          DEFAULT 0,
    total_dates     INT          DEFAULT 0,
    success_dates   INT          DEFAULT 0,
    avg_rating      DOUBLE       DEFAULT 5.0,
    total_rating    INT          DEFAULT 0,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_matchmaker (matchmaker_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Chat Messages Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS chat_messages (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sender_id       BIGINT UNSIGNED NOT NULL,
    receiver_id     BIGINT UNSIGNED NOT NULL,
    type            VARCHAR(20)  NOT NULL DEFAULT 'text',
    content         TEXT         NOT NULL,
    is_read         TINYINT(1)   DEFAULT 0,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_sender (sender_id),
    INDEX idx_receiver (receiver_id),
    INDEX idx_sender_receiver (sender_id, receiver_id),
    INDEX idx_is_read (is_read)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Chat Sessions Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS chat_sessions (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_a_id       BIGINT UNSIGNED NOT NULL,
    user_b_id       BIGINT UNSIGNED NOT NULL,
    last_message    VARCHAR(500),
    last_time       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    unread_a        INT          DEFAULT 0,
    unread_b        INT          DEFAULT 0,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_users (user_a_id, user_b_id),
    INDEX idx_user_a (user_a_id),
    INDEX idx_user_b (user_b_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Member Benefits Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS member_benefits (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    level           VARCHAR(20)  NOT NULL UNIQUE,
    daily_interact  INT          DEFAULT 5,
    unlimited_chat  TINYINT(1)   DEFAULT 0,
    view_who_liked  TINYINT(1)   DEFAULT 0,
    priority_match  TINYINT(1)   DEFAULT 0,
    advanced_filter TINYINT(1)   DEFAULT 0,
    video_chat      TINYINT(1)   DEFAULT 0,
    hide_online     TINYINT(1)   DEFAULT 0,
    no_ads          TINYINT(1)   DEFAULT 0,
    matchmaker_assist TINYINT(1) DEFAULT 0,
    price_per_month DOUBLE       DEFAULT 0,
    description     TEXT,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Member Orders Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS member_orders (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL,
    level           VARCHAR(20)  NOT NULL,
    months          INT          NOT NULL,
    amount          DOUBLE       NOT NULL,
    status          VARCHAR(20)  NOT NULL DEFAULT 'pending',
    paid_at         DATETIME     NULL,
    expire_at       DATETIME     NULL,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Interact Logs Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS interact_logs (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED NOT NULL,
    target_id       BIGINT UNSIGNED,
    action          VARCHAR(50),
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Sensitive Words Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS sensitive_words (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    word            VARCHAR(50)  NOT NULL UNIQUE,
    category        VARCHAR(20),
    INDEX idx_word (word)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- System Logs Table
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS system_logs (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         BIGINT UNSIGNED,
    module          VARCHAR(50),
    action          VARCHAR(50),
    ip              VARCHAR(50),
    detail          TEXT,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_module (module),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- Initial Data
-- -----------------------------------------------------------
INSERT IGNORE INTO sensitive_words (word, category) VALUES
('赌博', 'illegal'), ('色情', 'illegal'), ('贷款', 'illegal'),
('诈骗', 'illegal'), ('毒品', 'illegal'), ('枪支', 'illegal'),
('暴力', 'illegal'), ('恐怖', 'illegal'), ('传销', 'illegal');

INSERT IGNORE INTO member_benefits (level, daily_interact, unlimited_chat, view_who_liked, priority_match, advanced_filter, video_chat, hide_online, no_ads, matchmaker_assist, price_per_month, description) VALUES
('free',    5,   0, 0, 0, 0, 0, 0, 0, 0, 0,     '免费会员：每天5次互动'),
('silver',  20,  1, 1, 0, 0, 0, 0, 0, 0, 29.9,  '白银会员：无限聊天、查看谁喜欢我'),
('gold',    50,  1, 1, 1, 1, 1, 0, 0, 0, 99.9,  '黄金会员：优先匹配、高级筛选、视频聊天'),
('diamond', 999, 1, 1, 1, 1, 1, 1, 1, 1, 299.9, '钻石会员：全部特权+红娘服务');
