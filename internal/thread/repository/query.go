package repository

const (
	getThreadByID = `
SELECT title,
       author,
       forum,
       message,
       votes,
       slug,
       created
FROM threads
WHERE thread_id = $1;`

	checkExistThreadByID = `
SELECT EXISTS(SELECT 1 FROM threads WHERE thread_id = $1);`

	getThreadIDBySlug = `
SELECT thread_id
FROM threads
WHERE slug = $1;`

	createForumThread = `
INSERT INTO threads(title, author, forum, message, slug)
VALUES ($1, $2, $3, $4, $5) RETURNING thread_id, created;`

	updateThreadByID = `
UPDATE threads
SET title   = $2,
    message = $3
WHERE thread_id = $1
RETURNING author, forum, votes, slug, created;`

	insertPosts = "INSERT INTO posts(parent, author, message, forum, thread_id, created) VALUES "

	getPostsByFlatBegin = `
SELECT post_id, parent, author, message, is_edited, forum, created FROM posts WHERE thread_id = $1 `

	getPostsByTreeBegin = `
SELECT post_id, parent, author, message, is_edited, forum, created FROM posts WHERE thread_id = $1 `
)
