package dbpostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/arfan21/backend-test/config"
	"github.com/arfan21/backend-test/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxOpenConnection = 60
	connMaxLifetime   = 120
	maxIdleConns      = 30
	connMaxIdleTime   = 20
)

func NewPgx() (db *pgxpool.Pool, err error) {
	ctx := context.Background()
	fmt.Println("DB URL", config.GetConfig().Database.GetDatabaseURL())
	pgConfig, err := pgxpool.ParseConfig(config.GetConfig().Database.GetDatabaseURL())
	if err != nil {
		return nil, err
	}
	pgConfig.MaxConns = maxOpenConnection
	pgConfig.MaxConnIdleTime = connMaxIdleTime * time.Second
	pgConfig.MaxConnLifetime = connMaxLifetime * time.Second

	db, err = pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		err = fmt.Errorf("failed to connect to database: %w", err)
		return nil, err
	}

	if err = db.Ping(ctx); err != nil {
		err = fmt.Errorf("failed to ping database: %w", err)
		return nil, err
	}

	logger.Log(ctx).Info().Msg("dbpostgres: connection established")

	return db, nil
}

type Queryer interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}
