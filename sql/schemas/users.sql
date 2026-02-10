
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    username    VARCHAR(30) NOT NULL UNIQUE, -- Nome de exibição (@usuario)
    email       VARCHAR(255) NOT NULL UNIQUE, -- Para login e recuperação
    password    TEXT NOT NULL,               -- Hash da senha (nunca texto puro!)

    avatar_url  TEXT,    -- Caminho da foto de perfil
    bio         VARCHAR(160),                -- Pequena descrição

    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Índice para busca rápida por nome de usuário (ex: busca de perfis)
CREATE INDEX idx_users_username ON users(username);
