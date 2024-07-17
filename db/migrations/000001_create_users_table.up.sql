-- AuthService uchun users jadvali
CREATE TABLE if not exists users (
     UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    user_type VARCHAR(20) NOT NULL,
    bio TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE Table if not EXISTS refresh_token(
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    token TEXT UNIQUE NOT NULL,
    expires_at BIGINT,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);
