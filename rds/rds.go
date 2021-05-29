package rds

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgxpool"
	_ "github.com/lib/pq"
)

type rdsClient struct {
	connString string
	conn       QueryCloser
}

type RdsClient interface {
	Close()
	Querier(ctx context.Context, tx pgx.Tx) (Querier, error)
	TestConnection() error
}

type QueryCloser interface {
	Querier
	Close()
}

type Querier interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func NewRDSClient(connString string) RdsClient {
	return &rdsClient{
		connString: connString,
	}
}

func (db *rdsClient) Close() {
	if db.conn != nil {
		db.conn.Close()
		db.conn = nil
	}
}

func (db *rdsClient) Querier(ctx context.Context, tx pgx.Tx) (Querier, error) {
	if tx != nil {
		return tx, nil
	}

	if err := db.connect(ctx); err != nil {
		return nil, err
	}

	return db.conn, nil
}

func (db *rdsClient) connect(ctx context.Context) error {
	if db.conn != nil {
		return nil
	}

	var err error
	db.conn, err = pgxpool.Connect(ctx, db.connString)
	return err
}

func (db *rdsClient) TestConnection() error {
	dbConnection, err := sql.Open("postgres", db.connString)
	if err != nil {
		return err
	}
	defer dbConnection.Close()

	err = dbConnection.Ping()
	if err != nil {
		return err
	}

	return nil
}
