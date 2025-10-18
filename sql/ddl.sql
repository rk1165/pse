DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS requests;

CREATE VIRTUAL TABLE posts
    USING fts5
(
    title,
    url,
    content,
    tokenize = 'porter unicode61 remove_diacritics 1'
);

CREATE TABLE requests
(
    url     TEXT PRIMARY KEY ,
    title   TEXT,
    links   TEXT,
    content TEXT
);
