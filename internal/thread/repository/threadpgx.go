package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"db-performance-project/internal/models"
	"db-performance-project/internal/pkg"
	"db-performance-project/internal/pkg/sqltools"
)

type ThreadRepository interface {
	// Support
	GetThreadIDBySlug(ctx context.Context, thread *models.Thread) (*models.Thread, error)

	CreateThread(ctx context.Context, thread *models.Thread) (*models.Thread, error)
	CreatePostsByID(ctx context.Context, thread *models.Thread, posts []*models.Post) ([]*models.Post, error)
	GetDetailsThreadByID(ctx context.Context, thread *models.Thread) (*models.Thread, error)
	UpdateThreadByID(ctx context.Context, thread *models.Thread) (*models.Thread, error)

	// Posts
	GetPostsByIDFlat(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]*models.Post, error)
	GetPostsByIDTree(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]*models.Post, error)
	GetPostsByIDParentTree(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]*models.Post, error)
}

type threadPostgres struct {
	database *sqltools.Database
}

func NewThreadPostgres(database *sqltools.Database) ThreadRepository {
	return &threadPostgres{
		database,
	}
}

func (t threadPostgres) GetThreadIDBySlug(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
	res := &models.Thread{}

	errMain := sqltools.RunQuery(ctx, t.database.Connection, func(ctx context.Context, conn *sql.Conn) error {
		rowThread := conn.QueryRowContext(ctx, getThreadIDBySlug, thread.Slug)
		if errors.As(rowThread.Err(), sql.ErrNoRows) {
			return pkg.ErrSuchThreadNotFound
		}
		// if rowCounters.err() != nil {
		//	return errors.WithMessagef(pkg.ErrWorkDatabase,
		//		"Err: params input: query - [%s], values - [%s, %s, %s, %s]. Special error: [%s]",
		//		createUser, user.Nickname, user.FullName, user.About, user.Email, rowUser.Err())
		// }

		err := rowThread.Scan(&res.ID)
		if err != nil {
			return err
		}

		return nil
	})

	if errMain != nil {
		return nil, errMain
	}

	return res, nil
}
func (t threadPostgres) CreateThread(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
	errMain := sqltools.RunTxOnConn(ctx, pkg.TxInsertOptions, t.database.Connection, func(ctx context.Context, tx *sql.Tx) error {
		rowUser := tx.QueryRowContext(ctx, createForumThread, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug)
		if errors.Is(rowUser.Err(), sql.ErrTxDone) {
			return pkg.ErrSuchThreadExist
		}

		// else {
		//	return errors.WithMessagef(pkg.ErrWorkDatabase,
		//		"Err: params input: query - [%s], values - [%s, %s, %s, %s]. Special error: [%s]",
		//		createUser, user.Nickname, user.FullName, user.About, user.Email, rowUser.Err())
		// }

		err := rowUser.Scan(&thread.ID, &thread.Created)
		if err != nil {
			return err
		}

		return nil
	})

	if errMain != nil {
		return nil, errMain
	}

	return thread, nil
}

func (t threadPostgres) CreatePostsByID(ctx context.Context, thread *models.Thread, posts []*models.Post) ([]*models.Post, error) {
	// Defining sending parameters
	query := insertPosts

	countAttributes := strings.Count(query, ",") + 1

	pos := 0

	countInserts := len(posts)

	values := make([]interface{}, countInserts*countAttributes)

	insertTime := time.Now()

	for i := 0; i < len(posts); i++ {
		values[pos] = posts[i].Parent
		pos++
		values[pos] = posts[i].Author
		pos++
		values[pos] = posts[i].Message
		pos++
		values[pos] = posts[i].Forum
		pos++
		values[pos] = posts[i].Thread
		pos++
		values[pos] = insertTime.Format(time.RFC3339)
		pos++
	}

	insertStatement := sqltools.CreateFullQuery(query, countInserts, countAttributes)

	insertStatement += " RETURNING post_id;"

	rows, err := sqltools.InsertBatch(ctx, t.database.Connection, insertStatement, values)
	if err != nil {
		return nil, err
	}

	i := 0
	for rows.Next() {
		err = rows.Scan(&posts[i].ID)
		if err != nil {
			return nil, err
		}

		i++
	}

	return posts, nil
}

func (t threadPostgres) GetDetailsThreadByID(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
	errMain := sqltools.RunQuery(ctx, t.database.Connection, func(ctx context.Context, conn *sql.Conn) error {
		rowThread := conn.QueryRowContext(ctx, getThreadByID, thread.ID)
		if errors.As(rowThread.Err(), sql.ErrNoRows) {
			return pkg.ErrSuchThreadNotFound
		}
		// if rowCounters.err() != nil {
		//	return errors.WithMessagef(pkg.ErrWorkDatabase,
		//		"Err: params input: query - [%s], values - [%s, %s, %s, %s]. Special error: [%s]",
		//		createUser, user.Nickname, user.FullName, user.About, user.Email, rowUser.Err())
		// }

		err := rowThread.Scan(
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created)
		if err != nil {
			return err
		}

		return nil
	})

	if errMain != nil {
		return nil, errMain
	}

	return thread, nil
}

func (t threadPostgres) UpdateThreadByID(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
	errMain := sqltools.RunTxOnConn(ctx, pkg.TxInsertOptions, t.database.Connection, func(ctx context.Context, tx *sql.Tx) error {
		rowThread := tx.QueryRowContext(ctx, updateThreadByID, thread.ID, thread.Title, thread.Message)
		if errors.As(rowThread.Err(), sql.ErrNoRows) {
			return pkg.ErrSuchThreadNotFound
		}
		// if rowCounters.err() != nil {
		//	return errors.WithMessagef(pkg.ErrWorkDatabase,
		//		"Err: params input: query - [%s], values - [%s, %s, %s, %s]. Special error: [%s]",
		//		createUser, user.Nickname, user.FullName, user.About, user.Email, rowUser.Err())
		// }

		err := rowThread.Scan(
			&thread.Author,
			&thread.Forum,
			&thread.Votes,
			&thread.Slug,
			&thread.Created)
		if err != nil {
			return err
		}

		return nil
	})

	if errMain != nil {
		return nil, errMain
	}

	return thread, nil
}

func (t threadPostgres) GetPostsByIDFlat(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]*models.Post, error) {
	var rows *sql.Rows
	var err error

	query := getPostsByFlatBegin

	var values []interface{}

	switch {
	case params.Since != -1 && params.Desc:
		query += " AND post_id < $2"
	case params.Since != -1 && !params.Desc:
		query += " AND post_id > $2"
	case params.Since != -1:
		query += " AND post_id > $2"
	}

	switch {
	case params.Desc:
		query += " ORDER BY created DESC, post_id DESC"
	case !params.Desc:
		query += " ORDER BY created ASC, post_id"
	default:
		query += " ORDER BY created, post_id"
	}

	query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", params.Limit)

	if params.Since == -1 {
		values = []interface{}{thread.ID}
	} else {
		values = []interface{}{thread.ID, params.Since}
	}

	res := make([]*models.Post, 0)

	err = sqltools.RunQuery(ctx, t.database.Connection, func(ctx context.Context, conn *sql.Conn) error {
		rows, err = conn.QueryContext(ctx, query, values...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			post := &models.Post{}

			timeTmp := time.Time{}

			err = rows.Scan(
				&post.ID,
				&post.Parent,
				&post.Author,
				&post.Message,
				&post.IsEdited,
				&post.Forum,
				&timeTmp)
			if err != nil {
				return err
			}

			post.Thread = thread.ID

			post.Created = timeTmp.Format(time.RFC3339)

			res = append(res, post)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t threadPostgres) GetPostsByIDTree(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]*models.Post, error) {
	var rows *sql.Rows
	var err error

	query := getPostsByTreeBegin

	switch {
	case params.Since != -1 && params.Desc:
		query += " AND path < "
	case params.Since != -1 && !params.Desc:
		query += " AND path > "
	case params.Since != -1:
		query += " AND path > "
	}

	if params.Since != -1 {
		query += fmt.Sprintf(` (SELECT path FROM post WHERE post_id = %d) `, params.Since)
	}

	switch {
	case params.Desc:
		query += " ORDER BY path DESC"
	case !params.Desc:
		query += " ORDER BY path ASC, post_id"
	default:
		query += " ORDER BY path, post_id"
	}

	query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", params.Limit)

	res := make([]*models.Post, 0)

	err = sqltools.RunQuery(ctx, t.database.Connection, func(ctx context.Context, conn *sql.Conn) error {
		rows, err = conn.QueryContext(ctx, query, thread.ID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			post := &models.Post{}

			timeTmp := time.Time{}

			err = rows.Scan(
				&post.ID,
				&post.Parent,
				&post.Author,
				&post.Message,
				&post.IsEdited,
				&post.Forum,
				&timeTmp)
			if err != nil {
				return err
			}

			post.Thread = thread.ID

			post.Created = timeTmp.Format(time.RFC3339)

			res = append(res, post)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t threadPostgres) GetPostsByIDParentTree(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]*models.Post, error) {
	var rows *sql.Rows
	var err error

	query := ""

	var values []interface{}

	if params.Since == -1 {
		if params.Desc {
			query = `
					SELECT post_id, parent, author, message, is_edited, forum, created FROM posts
					WHERE path[1] IN (SELECT post_id FROM posts WHERE thread_id = $1 AND parent = 0 ORDER BY post_id DESC LIMIT $2)
					ORDER BY path[1] DESC, path ASC, post_id ASC;`
		} else {
			query = `
					SELECT post_id, parent, author, message, is_edited, forum, created FROM posts
					WHERE path[1] IN (SELECT post_id FROM posts WHERE thread_id = $1 AND parent = 0 ORDER BY post_id ASC LIMIT $2)
					ORDER BY path ASC, post_id ASC;`
		}

		values = []interface{}{thread.ID, params.Limit}
	} else {
		if params.Desc {
			query = `
					SELECT post_id, parent, author, message, is_edited, forum, created FROM posts
					WHERE path[1] IN (SELECT post_id FROM posts WHERE thread_id = $1 AND parent = 0 AND path[1] <
					(SELECT path[1] FROM posts WHERE post_id = $2) ORDER BY post_id DESC LIMIT $3)
					ORDER BY path[1] DESC, path ASC, post_id ASC;`
		} else {
			query = `
					SELECT post_id, parent, author, message, is_edited, forum, created FROM posts
					WHERE path[1] IN (SELECT post_id FROM posts WHERE thread_id = $1 AND parent = 0 AND path[1] >
					(SELECT path[1] FROM posts WHERE post_id = $2) ORDER BY post_id ASC LIMIT $3) 
					ORDER BY path ASC, post_id ASC;`
		}

		values = []interface{}{thread.ID, params.Since, params.Limit}
	}

	res := make([]*models.Post, 0)

	err = sqltools.RunQuery(ctx, t.database.Connection, func(ctx context.Context, conn *sql.Conn) error {
		rows, err = conn.QueryContext(ctx, query, values...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			post := &models.Post{}

			timeTmp := time.Time{}

			err = rows.Scan(
				&post.ID,
				&post.Parent,
				&post.Author,
				&post.Message,
				&post.IsEdited,
				&post.Forum,
				&timeTmp)
			if err != nil {
				return err
			}

			post.Thread = thread.ID

			post.Created = timeTmp.Format(time.RFC3339)

			res = append(res, post)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
