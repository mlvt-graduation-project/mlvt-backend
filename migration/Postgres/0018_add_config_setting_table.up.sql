CREATE TABLE feature_flags (
    id SERIAL PRIMARY KEY,
    flag_key TEXT UNIQUE NOT NULL,              -- khóa duy nhất cho từng flag (ví dụ: 'model_charge', 'token_limit')
    config_details JSONB NOT NULL,              -- chứa cấu hình động, có thể khác nhau mỗi dòng
    is_active BOOLEAN DEFAULT TRUE,             -- bật/tắt flag
    description TEXT,                           -- mô tả tùy chọn
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO feature_flags (flag_key, config_details, is_active, description)
VALUES (
  'model charge',
  '{
    "STT": { "default": 5 },
    "TTT": { "default": 2 },
    "TTS": { "default": 5 },
    "LS": { "default": 5 },
    "FP": { "default": 10 }
  }',
  false,
  'feature flag for charging the user if the pipeline function will cost token'
);
