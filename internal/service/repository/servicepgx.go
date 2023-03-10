package repository

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"db-performance-project/internal/models"
	"db-performance-project/internal/pkg"
	"db-performance-project/internal/pkg/sqltools"
)

type ServiceRepository interface {
	Clear(ctx context.Context) error
	GetStatus(ctx context.Context) (*models.StatusService, error)
}

type servicePostgres struct {
	database *sqltools.Database
}

func NewServicePostgres(database *sqltools.Database) ServiceRepository {
	return &servicePostgres{
		database,
	}
}

func (s servicePostgres) Clear(ctx context.Context) error {
	errMain := sqltools.RunTxOnConn(ctx, pkg.TxInsertOptions, s.database.Connection, func(ctx context.Context, tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, clearAllTables)
		if row.Err() != nil {
			return errors.WithMessagef(pkg.ErrWorkDatabase,
				"Err: params input: query - [%s]. Special error: [%s]",
				clearAllTables, row.Err())
		}

		return nil
	})

	return errMain
}

func (s servicePostgres) GetStatus(ctx context.Context) (*models.StatusService, error) {
	res := &models.StatusService{}

	errMain := sqltools.RunQuery(ctx, s.database.Connection, func(ctx context.Context, conn *sql.Conn) error {
		rowCounters := conn.QueryRowContext(ctx, getCountForumsPostsThreadsUsers)
		if rowCounters.Err() != nil {
			return errors.WithMessagef(pkg.ErrWorkDatabase,
				"Err: params input: query - [%s]. Special error: [%s]",
				getCountForumsPostsThreadsUsers, rowCounters.Err())
		}

		err := rowCounters.Scan(
			&res.Forum,
			&res.Post,
			&res.Thread,
			&res.User)
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
