package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"go.uber.org/multierr"
)

// SQLError holds an error related to a SQL.
type SQLError struct {
	query string
	args  []any
	err   error
}

// Error returns the error in string format.
func (e SQLError) Error() string {
	return fmt.Sprintf("sql::%s", e.err.Error())
}

// SQLRow is the supported function to for a returned row.
type SQLRow interface {
	Scan(dest ...any) error
	Columns() ([]string, error)
	ColumnTypes() ([]*sql.ColumnType, error)
}

// SQLScan scans the row.
type SQLScan interface {
	Scan(row SQLRow) error
}

// SQLQuerySession is a SQL query.
type SQLQuerySession struct {
	Query   string
	Args    []any
	Scanner SQLScan
}

// SQLQuery queries a SQL query.
func SQLQuery(ctx context.Context, conn *Conn, session SQLQuerySession) (results []SQLScan, errs error) {
	if session.Query == "" {
		return nil, nil
	}

	rows, err := conn.QueryContext(ctx, session.Query, session.Args...)
	if err != nil {
		return nil, SQLError{query: session.Query, args: session.Args, err: err}
	}
	defer rows.Close()

	for rows.Next() {
		if session.Scanner == nil {
			results = append(results, nil)
			continue
		}

		result := reflect.New(reflect.TypeOf(session.Scanner).Elem()).Interface().(SQLScan)
		if err := result.Scan(rows); err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		results = append(results, result)
	}
	return results, errs
}

// SQLExecSession is a SQL execution.
type SQLExecSession struct {
	Query    string
	Args     []any
	Callback func(sql.Result) error
}

// SQLExec exes a SQL.
func SQLExec(ctx context.Context, conn *Conn, session SQLExecSession) error {
	if session.Query == "" {
		return nil
	}
	if result, err := conn.ExecContext(ctx, session.Query, session.Args...); err != nil {
		return SQLError{query: session.Query, args: session.Args, err: err}
	} else if session.Callback != nil {
		if err := session.Callback(result); err != nil {
			return SQLError{query: session.Query, args: session.Args, err: err}
		}
	}
	return nil
}

// SQLBatchOption configures the SQLBatchExecSession.
type SQLBatchOption struct {
	BatchSize int `validate:"gt=0"`
}

var defaultBatchOption = SQLBatchOption{BatchSize: 10}

// SQLBatchExecSession creates the SQLExecSession in batch.
func SQLBatchExecSession[S any](data []S, builder func([]S) SQLExecSession, opts ...func(*SQLBatchOption)) (sessions []SQLExecSession) {
	opt := defaultBatchOption
	for _, fn := range opts {
		fn(&opt)
	}

	batches := make([][]S, 0, (len(data)+opt.BatchSize-1)/opt.BatchSize)
	for opt.BatchSize < len(data) {
		data, batches = data[opt.BatchSize:], append(batches, data[0:opt.BatchSize:opt.BatchSize])
	}
	batches = append(batches, data)

	for _, batch := range batches {
		sessions = append(sessions, builder(batch))
	}
	return sessions
}
