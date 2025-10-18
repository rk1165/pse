SELECT title, content
FROM posts;

SELECT count(*) as total
FROM posts;

SELECT title, content,  bm25(posts, 5, 2, 1) as rank
FROM posts
WHERE posts = 'university'
ORDER BY rank;

SELECT title, content
FROM posts
WHERE posts MATCH 'learn golang';

-- prefix search
SELECT title, content
FROM posts
WHERE posts MATCH 'learn*'
ORDER BY rank;

SELECT highlight(posts, 0, '<b>', '</b>') title,
       highlight(posts, 2, '<b>', '</b>') content
FROM posts
WHERE posts MATCH 'learn'
ORDER BY rank;

SELECT title, content
FROM posts
WHERE posts MATCH 'learn*'
ORDER BY rank
LIMIT 10;

SELECT title, content
FROM posts
WHERE posts MATCH 'github pages';

SELECT * FROM requests;

SELECT title, content
FROM posts
WHERE posts = 'university'
ORDER BY rank
LIMIT 10
OFFSET 20;
