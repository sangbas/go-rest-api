package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

// The operation constants
const (
	execOperation      = "repository.base.exec"
	fetchRowOperation  = "repository.base.fetch_row"
	fetchRowsOperation = "repository.base.fetch_rows"
)

// BaseRepository type
type BaseRepository struct {
	MasterDB *sqlx.DB
	SlaveDB  *sqlx.DB
}

func (r *BaseRepository) Exec(ctx context.Context, query string, args interface{}) (sql.Result, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, execOperation)
	defer span.Finish()

	var (
		res sql.Result
		err error
	)

	if r.MasterDB == nil {
		return res, errors.New("the master DB connection is nil")
	}

	res, err = r.MasterDB.NamedExecContext(ctx, query, args)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// FetchRow the fetch data row on Slave DB
func (r *BaseRepository) FetchRows(ctx context.Context, query string, resp interface{}, args ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, fetchRowOperation)
	defer span.Finish()

	if r.SlaveDB == nil {
		return errors.New("the slave DB connection is nil")
	}

	err := r.SlaveDB.Select(resp, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// FetchRow the fetch data row on Slave DB
func (r *BaseRepository) FetchRow(ctx context.Context, query string, resp interface{}, args ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, fetchRowOperation)
	defer span.Finish()

	if r.SlaveDB == nil {
		return errors.New("the slave DB connection is nil")
	}

	err := r.SlaveDB.Get(resp, query, args...)
	if err != nil {
		return err
	}

	return nil
}
