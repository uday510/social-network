DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_users_created_at') THEN
CREATE INDEX idx_users_created_at ON users(created_at);
END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_posts_user_id') THEN
CREATE INDEX idx_posts_user_id ON posts(user_id);
END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_posts_created_at') THEN
CREATE INDEX idx_posts_created_at ON posts(created_at);
END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_posts_tags_gin') THEN
CREATE INDEX idx_posts_tags_gin ON posts USING GIN (tags);
END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_comments_user_id') THEN
CREATE INDEX idx_comments_user_id ON comments(user_id);
END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_comments_post_id') THEN
CREATE INDEX idx_comments_post_id ON comments(post_id);
END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_comments_created_at') THEN
CREATE INDEX idx_comments_created_at ON comments(created_at);
END IF;
END$$;