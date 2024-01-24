package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/config"
	"log"
	"time"
)

func Open(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Name,
		time.Local.String(),
	)

	pgxpoolConfig, err := pgxpool.ParseConfig(psqlInfo)

	if err != nil {
		log.Fatalf("postgres client parse config error %v", err)
	}

	pgxpoolConfig.MaxConnLifetime = time.Duration(cfg.MaxConnLifeTime)
	pgxpoolConfig.MaxConnIdleTime = time.Second * 10
	pgxpoolConfig.MaxConns = cfg.MaxOpenConn

	pool, err := pgxpool.NewWithConfig(ctx, pgxpoolConfig)
	if err != nil {
		log.Fatalf("postgres client connect error %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("postgres client Ping error %v", err)
	}
	return pool, nil
}
