-- 001_init_schema.sql
-- 初始化数据库表结构

-- 词典表
CREATE TABLE IF NOT EXISTS dictionaries (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL DEFAULT 1,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    total_words INT DEFAULT 0,
    learned_words INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_dictionaries_user_id ON dictionaries(user_id);

-- 单词表
CREATE TABLE IF NOT EXISTS words (
    id BIGSERIAL PRIMARY KEY,
    dict_id BIGINT NOT NULL REFERENCES dictionaries(id) ON DELETE CASCADE,
    word VARCHAR(100) NOT NULL,
    phonetic VARCHAR(100),
    meaning JSONB NOT NULL DEFAULT '{}',
    example TEXT,
    audio_url VARCHAR(255),
    
    -- 记忆算法字段
    status VARCHAR(20) DEFAULT 'new',
    ef_factor NUMERIC(3,2) DEFAULT 2.50,
    interval INT DEFAULT 0,
    repetitions INT DEFAULT 0,
    next_review_date DATE,
    last_review_date DATE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(dict_id, word)
);

CREATE INDEX IF NOT EXISTS idx_words_dict_id ON words(dict_id);
CREATE INDEX IF NOT EXISTS idx_words_status ON words(status);
CREATE INDEX IF NOT EXISTS idx_words_next_review ON words(next_review_date);

-- 学习记录表
CREATE TABLE IF NOT EXISTS learn_records (
    id BIGSERIAL PRIMARY KEY,
    word_id BIGINT NOT NULL REFERENCES words(id) ON DELETE CASCADE,
    quality INT NOT NULL CHECK (quality >= 0 AND quality <= 5),
    time_spent INT,
    ef_factor_before NUMERIC(3,2),
    ef_factor_after NUMERIC(3,2),
    interval_before INT,
    interval_after INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_learn_records_word_id ON learn_records(word_id);
CREATE INDEX IF NOT EXISTS idx_learn_records_created_at ON learn_records(created_at DESC);

-- 上传任务表（用于跟踪异步任务进度）
CREATE TABLE IF NOT EXISTS upload_tasks (
    id VARCHAR(64) PRIMARY KEY,
    dict_id BIGINT REFERENCES dictionaries(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'pending',
    total_words INT DEFAULT 0,
    processed_words INT DEFAULT 0,
    failed_words JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_upload_tasks_status ON upload_tasks(status);
