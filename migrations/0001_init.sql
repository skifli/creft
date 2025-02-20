-- Admins per Guild (Many-to-Many)
CREATE TABLE IF NOT EXISTS admins (
    guild_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    PRIMARY KEY (guild_id, user_id)
);

-- Counting Channels
CREATE TABLE IF NOT EXISTS counting (
    channel_id BIGINT PRIMARY KEY,
    guild_id BIGINT NOT NULL DEFAULT 0,
    count BIGINT NOT NULL DEFAULT 0,
    count_max BIGINT NOT NULL DEFAULT 0,
    last_count_message_edited BOOLEAN NOT NULL DEFAULT FALSE,
    last_count_message_id BIGINT NOT NULL DEFAULT 0,
    last_count_user_id BIGINT NOT NULL DEFAULT 0,
    resets_count BIGINT NOT NULL DEFAULT 0
);

-- Counting Stats (Per Channel, Per User, Linked to `counting` and `guild`)
CREATE TABLE IF NOT EXISTS counting_stats (
    channel_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    guild_id BIGINT NOT NULL,
    deleted_count_message BIGINT NOT NULL DEFAULT 0,
    edited_count_message BIGINT NOT NULL DEFAULT 0,
    correct BIGINT NOT NULL DEFAULT 0,
    incorrect BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (channel_id, user_id),  -- Composite key (channel + user)
    FOREIGN KEY (channel_id) REFERENCES counting (channel_id) ON DELETE CASCADE
);
