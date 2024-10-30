-- Filename: migrations/000002_add_comments_likes_column.down.sql
ALTER TABLE comments
DROP COLUMN IF EXISTS likes;
