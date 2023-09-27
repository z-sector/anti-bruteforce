package pgdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type PostgresDB struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

func MustPostgresDB(url string, opts ...Option) *PostgresDB {
	p, err := NewPostgresDB(url, opts...)
	if err != nil {
		panic(err)
	}
	return p
}

func NewPostgresDB(url string, opts ...Option) (*PostgresDB, error) {
	pg := &PostgresDB{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("pgxpoll parse config: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			if err = pg.Pool.Ping(context.Background()); err == nil {
				break
			}
		}

		log.Printf("postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (p *PostgresDB) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
