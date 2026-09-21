-- wishwall 数据库初始化脚本（GORM AutoMigrate 为实际执行方，本脚本供人工参考/审计）
CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(120) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname      VARCHAR(50)  NOT NULL,
    avatar        VARCHAR(255) DEFAULT '',
    bio           VARCHAR(500) DEFAULT '',
    role          VARCHAR(20)  NOT NULL DEFAULT 'user',
    status        VARCHAR(20)  NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS wishes (
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT NOT NULL REFERENCES users(id),
    title             VARCHAR(100) NOT NULL,
    content           TEXT NOT NULL,
    image_urls        TEXT NOT NULL DEFAULT '[]',
    category          VARCHAR(30) NOT NULL DEFAULT 'other',
    visibility        VARCHAR(20) NOT NULL DEFAULT 'public',
    difficulty        VARCHAR(20) NOT NULL DEFAULT 'medium',
    expected_deadline TIMESTAMPTZ,
    status            VARCHAR(20) NOT NULL DEFAULT 'pending',
    likes_count       INT NOT NULL DEFAULT 0,
    completion_note   TEXT DEFAULT '',
    is_anonymous      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_wishes_user ON wishes(user_id);
CREATE INDEX IF NOT EXISTS idx_wishes_status ON wishes(status);
CREATE INDEX IF NOT EXISTS idx_wishes_visibility ON wishes(visibility);
CREATE INDEX IF NOT EXISTS idx_wishes_category ON wishes(category);

CREATE TABLE IF NOT EXISTS wish_claims (
    id              BIGSERIAL PRIMARY KEY,
    wish_id         BIGINT NOT NULL UNIQUE REFERENCES wishes(id),
    user_id         BIGINT NOT NULL REFERENCES users(id),
    progress        INT NOT NULL DEFAULT 0,
    latest_note     TEXT DEFAULT '',
    status          VARCHAR(20) NOT NULL DEFAULT 'claimed',
    milestone_count INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_claims_user ON wish_claims(user_id);

CREATE TABLE IF NOT EXISTS blessings (
    id             BIGSERIAL PRIMARY KEY,
    wish_id        BIGINT NOT NULL REFERENCES wishes(id),
    user_id        BIGINT NOT NULL REFERENCES users(id),
    content        TEXT NOT NULL,
    gift_emoji     VARCHAR(50) DEFAULT '',
    is_celebrating BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_blessings_wish ON blessings(wish_id);

CREATE TABLE IF NOT EXISTS time_capsules (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    title       VARCHAR(100) NOT NULL,
    content     TEXT NOT NULL,
    image_urls  TEXT NOT NULL DEFAULT '[]',
    audio_url   VARCHAR(255) DEFAULT '',
    unlock_at   TIMESTAMPTZ NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'locked',
    unlocked_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_capsules_user ON time_capsules(user_id);

CREATE TABLE IF NOT EXISTS badges (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    type        VARCHAR(50) NOT NULL,
    title       VARCHAR(100) NOT NULL,
    description VARCHAR(255) DEFAULT '',
    icon        VARCHAR(50) DEFAULT '',
    earned_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, type)
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL DEFAULT 0,
    action      VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) DEFAULT '',
    entity_id   VARCHAR(64) DEFAULT '',
    detail      TEXT DEFAULT '',
    ip          VARCHAR(64) DEFAULT '',
    request_id  VARCHAR(64) DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_audit_user ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_request ON audit_logs(request_id);
