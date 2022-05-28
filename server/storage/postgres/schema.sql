CREATE TABLE IF NOT EXISTS users (
    id uuid UNIQUE PRIMARY KEY,
    email text UNIQUE NOT NULL,
    password_hash text NOT NULL,
    data_version integer NOT NULL DEFAULT 0,
    created_at timestamp,
    deleted_at timestamp
);

CREATE TABLE IF NOT EXISTS sessions (
    id uuid UNIQUE NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL,
    refresh_token text,
    login_at timestamp,
    logout_at timestamp,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS passwords (
    id uuid UNIQUE NOT NULL,
    user_id uuid NOT NULL,
    password TEXT NOT NULL,
    version integer NOT NULL DEFAULT 0,
    meta TEXT,
    created_at timestamp,
    deleted_at timestamp,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE TABLE IF NOT EXISTS blobs (
    id uuid UNIQUE NOT NULL,
    user_id uuid NOT NULL,
    blob BYTEA,
    version integer NOT NULL DEFAULT 0,
    meta TEXT,
    created_at timestamp,
    deleted_at timestamp,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE TABLE IF NOT EXISTS texts (
    id uuid UNIQUE NOT NULL,
    user_id uuid NOT NULL,
    text_string text,
    version integer NOT NULL DEFAULT 0,
    meta TEXT,
    created_at timestamp,
    deleted_at timestamp,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE TABLE IF NOT EXISTS cards (
    id uuid UNIQUE NOT NULL,
    user_id uuid NOT NULL,
    card_number text,
    cardholder_name text,
    expiration_date text,
    cvc integer,
    version integer NOT NULL DEFAULT 0,
    meta TEXT,
    created_at timestamp,
    deleted_at timestamp,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);