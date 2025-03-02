CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE secure_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    data_type TEXT NOT NULL CHECK (data_type IN ('password', 'text', 'binary', 'card')), 
    data_encrypted TEXT NOT NULL,
    meta_info TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
