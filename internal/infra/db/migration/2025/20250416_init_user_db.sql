-- Create AccountStatus enum type
CREATE TYPE account_status AS ENUM (
    'active',
    'inactive',
    'pending',
    'suspended',
    'banned',
    'deleted',
    'archived'
);

-- Create Role enum type
CREATE TYPE role AS ENUM (
    'basic',
    'premium',
    'cj',
    'admin'
);

-- Create the users table with the appropriate enum types for `status` and `role`
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    status account_status NOT NULL,  -- Use the AccountStatus enum
    role role NOT NULL,              -- Use the Role enum
    avatar VARCHAR(255),
    avatar_folder VARCHAR(255),
    wallet_balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Optionally, you can add an index to improve search performance on frequently queried fields
CREATE INDEX idx_users_status ON users (status);
CREATE INDEX idx_users_role ON users (role);
