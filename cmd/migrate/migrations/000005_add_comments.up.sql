CREATE TABLE IF NOT EXISTS comments (
                                        id bigserial PRIMARY KEY,
                                        post_id bigint NOT NULL,
                                        user_id bigint NOT NULL,
                                        content TEXT NOT NULL,
                                        created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),

                                        FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
                                        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);