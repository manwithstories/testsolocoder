CREATE DATABASE IF NOT EXISTS gym_management DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE gym_management;

CREATE TABLE IF NOT EXISTS members (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(100),
    gender VARCHAR(10),
    birthday DATE,
    address VARCHAR(255),
    profile_photo VARCHAR(255),
    password VARCHAR(255) NOT NULL,
    status TINYINT DEFAULT 1 COMMENT '1:正常 2:冻结',
    membership_id BIGINT UNSIGNED,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    UNIQUE KEY idx_phone (phone),
    UNIQUE KEY idx_email (email),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS memberships (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT UNSIGNED NOT NULL,
    type VARCHAR(20) NOT NULL COMMENT 'monthly:月卡 quarter:季卡 yearly:年卡',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    status TINYINT DEFAULT 1 COMMENT '1:有效 2:已过期 3:已冻结',
    auto_renew BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_member_id (member_id),
    INDEX idx_status (status),
    INDEX idx_end_date (end_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS coaches (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    specialty VARCHAR(255),
    description TEXT,
    photo VARCHAR(255),
    status TINYINT DEFAULT 1 COMMENT '1:在职 2:离职',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    UNIQUE KEY idx_phone (phone),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS courses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    coach_id BIGINT UNSIGNED NOT NULL,
    capacity INT NOT NULL,
    duration INT NOT NULL COMMENT '分钟',
    type VARCHAR(20) NOT NULL DEFAULT 'single' COMMENT 'single:单次 weekly:周课 monthly:月课',
    weekdays VARCHAR(20) COMMENT '例如: 1,3,5 表示周一、三、五',
    start_date DATE NOT NULL,
    end_date DATE,
    start_time VARCHAR(10) NOT NULL,
    location VARCHAR(100),
    status TINYINT DEFAULT 1 COMMENT '1:正常 2:取消 3:结束',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_coach_id (coach_id),
    INDEX idx_status (status),
    INDEX idx_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS course_schedules (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    course_id BIGINT UNSIGNED NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    capacity INT NOT NULL,
    booked_count INT DEFAULT 0,
    status TINYINT DEFAULT 1 COMMENT '1:可预约 2:已满 3:已取消 4:已结束',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_course_id (course_id),
    INDEX idx_start_time (start_time),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS bookings (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT UNSIGNED NOT NULL,
    schedule_id BIGINT UNSIGNED NOT NULL,
    status TINYINT DEFAULT 1 COMMENT '1:已预约 2:已取消 3:已签到 4:未到场',
    booking_time DATETIME NOT NULL,
    cancel_time DATETIME,
    check_in_id BIGINT UNSIGNED,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    UNIQUE KEY idx_member_schedule (member_id, schedule_id),
    INDEX idx_schedule_id (schedule_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS waitlists (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT UNSIGNED NOT NULL,
    schedule_id BIGINT UNSIGNED NOT NULL,
    position INT NOT NULL,
    notified BOOLEAN DEFAULT FALSE,
    status TINYINT DEFAULT 1 COMMENT '1:等待中 2:已转预约 3:已取消',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    UNIQUE KEY idx_waitlist_member_schedule (member_id, schedule_id),
    INDEX idx_schedule_id (schedule_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS check_ins (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT UNSIGNED NOT NULL,
    schedule_id BIGINT UNSIGNED,
    check_in_time DATETIME NOT NULL,
    check_type TINYINT DEFAULT 1 COMMENT '1:正常签到 2:补签',
    remark VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_member_id (member_id),
    INDEX idx_schedule_id (schedule_id),
    INDEX idx_check_in_time (check_in_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS reminders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT UNSIGNED NOT NULL,
    schedule_id BIGINT UNSIGNED,
    type TINYINT NOT NULL COMMENT '1:会员卡到期提醒 2:课程开始前1天提醒 3:课程开始前2小时提醒',
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    schedule_time DATETIME NOT NULL,
    sent_time DATETIME,
    status TINYINT DEFAULT 1 COMMENT '1:待发送 2:已发送 3:发送失败',
    retry_count INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_member_id (member_id),
    INDEX idx_schedule_time (schedule_time),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT UNSIGNED,
    operator_id BIGINT UNSIGNED,
    operator_type VARCHAR(20),
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id BIGINT UNSIGNED,
    detail TEXT,
    ip_address VARCHAR(50),
    user_agent VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_member_id (member_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO coaches (name, phone, specialty, description) VALUES 
('张教练', '13800138001', '瑜伽、普拉提', '10年瑜伽教学经验，国际认证教练'),
('李教练', '13800138002', '力量训练、HIIT', '前国家队体能教练'),
('王教练', '13800138003', '动感单车、有氧操', '健身教练国家职业资格认证');

INSERT INTO members (name, phone, email, password, status) VALUES 
('测试会员', '13800138000', 'test@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1);
