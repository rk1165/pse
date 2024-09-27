drop table if exists posts;
drop table if exists requests;

create virtual table posts
    using fts5
(
    title,
    url,
    content,
    tokenize = 'porter unicode61 remove_diacritics 1'
);

create table requests
(
    url     text primary key,
    title   text,
    links   text,
    content text
);


-- delete
-- from posts;
--
-- delete
-- from requests;
