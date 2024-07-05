CREATE TABLE IF NOT EXISTS rooms (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,

    name VARCHAR(255) NOT NULL,
    password_protected BOOL NOT NULL,
    password VARCHAR(255),
    password_salt VARCHAR(255)
);