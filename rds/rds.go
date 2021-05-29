package rds

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgxpool"
)

type rdsClient struct {
	secretName string
	conn       QueryCloser
}

type RdsClient interface {
	Close()
	Querier(ctx context.Context, tx pgx.Tx) (Querier, error)
}

type QueryCloser interface {
	Querier
	Close()
}

type Querier interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

func NewRDSClient(secretName string) RdsClient {
	return &rdsClient{
		secretName: secretName,
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
	db.conn, err = pgxpool.Connect(ctx, db.secretName)
	return err
}
