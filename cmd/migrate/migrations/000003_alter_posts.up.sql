DO $$
BEGIN
ALTER TABLE posts
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id);
EXCEPTION
    WHEN duplicate_object THEN
        RAISE NOTICE 'Constraint fk_user already exists.';
END;
$$;