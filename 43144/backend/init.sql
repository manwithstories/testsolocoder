CREATE DATABASE pet_adoption;

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL DEFAULT 'adopter',
    rescue_id BIGINT REFERENCES rescue_stations(id),
    address VARCHAR(500),
    is_verified BOOLEAN DEFAULT FALSE,
    avatar VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE rescue_stations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(500),
    contact_person VARCHAR(100),
    contact_phone VARCHAR(20),
    contact_email VARCHAR(255),
    license_number VARCHAR(100),
    license_file VARCHAR(500),
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    verified_by BIGINT REFERENCES users(id),
    verified_at TIMESTAMP,
    reject_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE pets (
    id BIGSERIAL PRIMARY KEY,
    archive_number VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    species VARCHAR(50) NOT NULL,
    breed VARCHAR(100),
    age VARCHAR(50),
    gender VARCHAR(20),
    weight DECIMAL(10,2),
    color VARCHAR(50),
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'adoptable',
    photos TEXT,
    videos TEXT,
    health_status VARCHAR(255),
    vaccinated BOOLEAN DEFAULT FALSE,
    neutered BOOLEAN DEFAULT FALSE,
    rescue_id BIGINT NOT NULL REFERENCES rescue_stations(id),
    adopter_id BIGINT REFERENCES users(id),
    found_location VARCHAR(255),
    found_date DATE,
    adopted_date DATE,
    medical_history TEXT,
    personality TEXT,
    special_needs TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE adoption_applications (
    id BIGSERIAL PRIMARY KEY,
    pet_id BIGINT NOT NULL REFERENCES pets(id),
    adopter_id BIGINT NOT NULL REFERENCES users(id),
    rescue_id BIGINT NOT NULL REFERENCES rescue_stations(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    reason TEXT,
    living_situation TEXT,
    pet_experience TEXT,
    family_members INTEGER,
    has_children BOOLEAN DEFAULT FALSE,
    has_other_pets BOOLEAN DEFAULT FALSE,
    other_pets_desc TEXT,
    housing_type VARCHAR(50),
    income_level VARCHAR(50),
    can_afford_vet BOOLEAN DEFAULT TRUE,
    agree_to_visit BOOLEAN DEFAULT TRUE,
    reviewed_by BIGINT REFERENCES users(id),
    reviewed_at TIMESTAMP,
    reject_reason TEXT,
    signed_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE adoption_agreements (
    id BIGSERIAL PRIMARY KEY,
    application_id BIGINT UNIQUE NOT NULL REFERENCES adoption_applications(id),
    adopter_sign BOOLEAN DEFAULT FALSE,
    adopter_signed_at TIMESTAMP,
    rescue_sign BOOLEAN DEFAULT FALSE,
    rescue_signed_at TIMESTAMP,
    agreement_terms TEXT NOT NULL,
    agreement_file VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE follow_up_records (
    id BIGSERIAL PRIMARY KEY,
    application_id BIGINT NOT NULL REFERENCES adoption_applications(id),
    pet_id BIGINT NOT NULL REFERENCES pets(id),
    adopter_id BIGINT NOT NULL REFERENCES users(id),
    rescue_id BIGINT NOT NULL REFERENCES rescue_stations(id),
    follow_up_date DATE,
    health_status TEXT,
    living_condition TEXT,
    notes TEXT,
    photo_evidence TEXT,
    recorded_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE health_records (
    id BIGSERIAL PRIMARY KEY,
    pet_id BIGINT NOT NULL REFERENCES pets(id),
    record_type VARCHAR(20) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    vaccine_name VARCHAR(100),
    record_date DATE,
    next_date DATE,
    weight DECIMAL(10,2),
    temperature DECIMAL(5,2),
    vet_name VARCHAR(100),
    hospital VARCHAR(200),
    report_file VARCHAR(500),
    notes TEXT,
    recorded_by BIGINT NOT NULL REFERENCES users(id),
    rescue_id BIGINT NOT NULL REFERENCES rescue_stations(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE health_reminders (
    id BIGSERIAL PRIMARY KEY,
    pet_id BIGINT NOT NULL REFERENCES pets(id),
    record_id BIGINT REFERENCES health_records(id),
    title VARCHAR(255) NOT NULL,
    reminder_date DATE NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE appointments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    pet_id BIGINT NOT NULL REFERENCES pets(id),
    rescue_id BIGINT NOT NULL REFERENCES rescue_stations(id),
    appointment_type VARCHAR(20) NOT NULL,
    appointment_date DATE NOT NULL,
    start_time VARCHAR(10) NOT NULL,
    end_time VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    location VARCHAR(255),
    notes TEXT,
    cancel_reason TEXT,
    original_id BIGINT REFERENCES appointments(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_pets_status ON pets(status);
CREATE INDEX idx_pets_rescue_id ON pets(rescue_id);
CREATE INDEX idx_pets_archive_number ON pets(archive_number);
CREATE INDEX idx_adoption_applications_status ON adoption_applications(status);
CREATE INDEX idx_adoption_applications_pet_id ON adoption_applications(pet_id);
CREATE INDEX idx_adoption_applications_adopter_id ON adoption_applications(adopter_id);
CREATE INDEX idx_adoption_applications_rescue_id ON adoption_applications(rescue_id);
CREATE INDEX idx_health_records_pet_id ON health_records(pet_id);
CREATE INDEX idx_health_records_record_type ON health_records(record_type);
CREATE INDEX idx_appointments_user_id ON appointments(user_id);
CREATE INDEX idx_appointments_pet_id ON appointments(pet_id);
CREATE INDEX idx_appointments_rescue_id ON appointments(rescue_id);
CREATE INDEX idx_appointments_status ON appointments(status);
CREATE INDEX idx_follow_up_records_pet_id ON follow_up_records(pet_id);
CREATE INDEX idx_health_reminders_pet_id ON health_reminders(pet_id);
