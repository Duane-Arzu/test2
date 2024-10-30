-- Filename: migrations/000002_add_comments_likes_column.up.sql
ALTER TABLE comments
ADD COLUMN likes integer NOT NULL DEFAULT 0;
