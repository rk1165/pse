select title, content
from posts;

select count(*) as total
from posts;

select title, content,  bm25(posts, 5, 2, 1) as rank
from posts
where posts = 'university'
order by rank;

select title, content
from posts
where posts match 'learn golang';

-- prefix search
select title, content
from posts
where posts match 'learn*'
order by rank;

select highlight(posts, 0, '<b>', '</b>') title,
       highlight(posts, 2, '<b>', '</b>') content
from posts
where posts match 'learn'
order by rank;

select title, content
from posts
where posts match 'learn*'
order by rank
limit 10;

select title, content
from posts
where posts match 'github pages';

select * from requests;

select title, content
from posts
where posts = 'university'
order by rank
limit 10
OFFSET 20;