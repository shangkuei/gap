package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"go.uber.org/multierr"
)

type txCtx struct{}

type txKey struct {
	tx   *sql.Tx
	uuid uuid.UUID
}

// TxError holds an error related to the SQL transaction.
type TxError struct {
	key *txKey
	err error
}

// Error returns the error in string format.
func (e TxError) Error() string {
	if e.key != nil {
		return fmt.Sprintf("tx(%s)::%s", e.key.uuid.String(), e.err.Error())
	}
	return fmt.Sprintf("tx::%s", e.err.Error())
}

// WithTx creates a new context with a transaction.
func WithTx(ctx context.Context, db *sql.DB) (context.Context, func(error) error, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return ctx, nil, TxError{err: err}
	}

	key := txKey{uuid: uuid.New(), tx: tx}
	ctx = context.WithValue(ctx, txCtx{}, &key)
	cancel := func(err error) error {
		if err != nil {
			slog.Debug("rollback transaction", "id", key.uuid)
			if txErr := tx.Rollback(); txErr != nil {
				return multierr.Append(err, TxError{key: &key, err: err})
			}
			return err
		}

		slog.Debug("commit transaction", "id", key.uuid)
		if err := tx.Commit(); err != nil {
			return TxError{key: &key, err: err}
		}
		return nil
	}
	return ctx, cancel, nil
}

// NewConn gets the Conn with the transaction in context or from a SQL connection.
func NewConn(ctx context.Context, db SQL, opts ...func(*Conn)) (conn *Conn) {
	conn = &Conn{SQL: db, Charset: "utf8mb4"}
	if key, ok := ctx.Value(txCtx{}).(*txKey); ok {
		slog.Debug("use tx", "id", key.uuid)
		conn.SQL = key.tx
	}

	for _, opt := range opts {
		opt(conn)
	}
	return conn
}

// Conn implements DB helpers.
type Conn struct {
	SQL
	Charset string
}

// SQL holds the basic SQL interface.
type SQL interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// WithCharset sets the default charset
func WithCharset(charset string) func(*Conn) {
	return func(conn *Conn) {
		conn.Charset = charset
	}
}
