DROP TABLE IF EXISTS posts;

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT[] NOT NULL
);

INSERT INTO posts (id, title, body) VALUES
(1, 'Post 1', ARRAY['Paragraph A', 'Paragraph B']),
(2, 'Post 2', ARRAY['Introduction', 'Conclusion']),
(3, 'Post 3', ARRAY['Single paragraph body.']);
