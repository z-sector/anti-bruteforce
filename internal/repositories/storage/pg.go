package storage

import (
	"context"
	"fmt"
	"net"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"anti_bruteforce/pkg/pgdb"
)

const (
	tableNameWhiteList = "white_list"
	tableNameBlackList = "black_list"
	columnIPAddress    = "ip_address"
)

type PgStorage struct {
	pgClient *pgdb.PostgresDB
	log      zerolog.Logger
}

func NewPgStorage(log zerolog.Logger, pgClient *pgdb.PostgresDB) *PgStorage {
	return &PgStorage{pgClient: pgClient, log: log}
}

func (p *PgStorage) CreateSubnetInBlackList(ctx context.Context, ipNet *net.IPNet) error {
	sql, args, err := p.pgClient.Builder.
		Insert(tableNameBlackList).Columns(columnIPAddress).
		Values(ipNet.String()).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgStorage - CreateSubnetInBlackList - pgClient.Builder: %w", err)
	}

	if _, err := p.pgClient.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("PgStorage - CreateSubnetInBlackList - pgClient.Pool.Exec: %w", err)
	}
	return nil
}

func (p *PgStorage) ExistsSubnetInBlackList(ctx context.Context, ipNet *net.IPNet) (bool, error) {
	sql, args, err := p.pgClient.Builder.
		Select("1").
		Prefix("SELECT EXISTS (").
		From(tableNameBlackList).
		Where(squirrel.Eq{columnIPAddress: ipNet.String()}).
		Suffix(")").
		ToSql()
	if err != nil {
		return false, fmt.Errorf("PgStorage - ExistsSubnetInBlackList - pgClient.Builder: %w", err)
	}

	var exists bool
	if err = p.pgClient.Pool.QueryRow(ctx, sql, args...).Scan(&exists); err != nil {
		return exists, fmt.Errorf("PgStorage - ExistsSubnetInBlackList - pgClient.Pool.QueryRow: %w", err)
	}
	return exists, nil
}

func (p *PgStorage) DeleteSubnetInBlackList(ctx context.Context, ipNet *net.IPNet) error {
	sql, args, err := p.pgClient.Builder.
		Delete(tableNameBlackList).
		Where(squirrel.Eq{columnIPAddress: ipNet.String()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgStorage - DeleteSubnetInBlackList - pgClient.Builder: %w", err)
	}

	if _, err := p.pgClient.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("PgStorage - DeleteSubnetInBlackList - pgClient.Pool.Exec: %w", err)
	}

	return nil
}

func (p *PgStorage) CreateSubnetInWhiteList(ctx context.Context, ipNet *net.IPNet) error {
	sql, args, err := p.pgClient.Builder.
		Insert(tableNameWhiteList).Columns(columnIPAddress).
		Values(ipNet.String()).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgStorage - CreateSubnetInWhiteList - pgClient.Builder: %w", err)
	}

	if _, err := p.pgClient.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("PgStorage - CreateSubnetInWhiteList - pgClient.Pool.Exec: %w", err)
	}
	return nil
}

func (p *PgStorage) ExistsSubnetInWhiteList(ctx context.Context, ipNet *net.IPNet) (bool, error) {
	sql, args, err := p.pgClient.Builder.
		Select("1").
		Prefix("SELECT EXISTS (").
		From(tableNameWhiteList).
		Where(squirrel.Eq{columnIPAddress: ipNet.String()}).
		Suffix(")").
		ToSql()
	if err != nil {
		return false, fmt.Errorf("PgStorage - ExistsSubnetInWhiteList - pgClient.Builder: %w", err)
	}

	var exists bool
	if err := p.pgClient.Pool.QueryRow(ctx, sql, args...).Scan(&exists); err != nil {
		return exists, fmt.Errorf("PgStorage - ExistsSubnetInWhiteList - pgClient.Pool.QueryRow: %w", err)
	}
	return exists, nil
}

func (p *PgStorage) DeleteSubnetInWhiteList(ctx context.Context, ipNet *net.IPNet) error {
	sql, args, err := p.pgClient.Builder.
		Delete(tableNameWhiteList).
		Where(squirrel.Eq{columnIPAddress: ipNet.String()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgStorage - DeleteSubnetInWhiteList - pgClient.Builder: %w", err)
	}

	if _, err := p.pgClient.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("PgStorage - DeleteSubnetInWhiteList - pgClient.Pool.Exec: %w", err)
	}

	return nil
}

func (p *PgStorage) ClearLists(ctx context.Context) error {
	sqlWhiteList, _, err := p.pgClient.Builder.
		Delete(tableNameWhiteList).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgStorage - ClearLists - pgClient.Builder: %w", err)
	}
	sqlBlackList, _, err := p.pgClient.Builder.
		Delete(tableNameBlackList).
		ToSql()
	if err != nil {
		return fmt.Errorf("PgStorage - ClearLists - pgClient.Builder: %w", err)
	}

	err = pgx.BeginFunc(ctx, p.pgClient.Pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, sqlWhiteList); err != nil {
			return fmt.Errorf("PgStorage - ClearLists - tx.Exec: %w", err)
		}
		if _, err := tx.Exec(ctx, sqlBlackList); err != nil {
			return fmt.Errorf("PgStorage - ClearLists - tx.Exec: %w", err)
		}
		return nil
	})

	return err
}

func (p *PgStorage) CheckInWhiteList(ctx context.Context, ip net.IP) (bool, error) {
	sql, args, err := p.pgClient.Builder.
		Select("1").
		Prefix("SELECT EXISTS (").
		From(tableNameWhiteList).
		Where(fmt.Sprintf("%s >>= $1", columnIPAddress), ip.String()).
		Suffix(")").
		ToSql()
	if err != nil {
		return false, fmt.Errorf("PgStorage - CheckInWhiteList - pgClient.Builder: %w", err)
	}

	var exists bool
	if err := p.pgClient.Pool.QueryRow(ctx, sql, args...).Scan(&exists); err != nil {
		return exists, fmt.Errorf("PgStorage - CheckInWhiteList - pgClient.Pool.QueryRow: %w", err)
	}
	return exists, nil
}

func (p *PgStorage) CheckInBlackList(ctx context.Context, ip net.IP) (bool, error) {
	sql, args, err := p.pgClient.Builder.
		Select("1").
		Prefix("SELECT EXISTS (").
		From(tableNameBlackList).
		Where(fmt.Sprintf("%s >>= $1", columnIPAddress), ip.String()).
		Suffix(")").
		ToSql()
	if err != nil {
		return false, fmt.Errorf("PgStorage - CheckInBlackList - pgClient.Builder: %w", err)
	}

	var exists bool
	if err := p.pgClient.Pool.QueryRow(ctx, sql, args...).Scan(&exists); err != nil {
		return exists, fmt.Errorf("PgStorage - CheckInBlackList - pgClient.Pool.QueryRow: %w", err)
	}
	return exists, nil
}
