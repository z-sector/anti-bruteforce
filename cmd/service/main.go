package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"anti_bruteforce/config"
	"anti_bruteforce/internal"
	"anti_bruteforce/internal/delivery/grpc"
	"anti_bruteforce/internal/repositories/bucket"
	"anti_bruteforce/internal/repositories/storage"
	"anti_bruteforce/internal/usecase"
	"anti_bruteforce/pkg/logger"
	"anti_bruteforce/pkg/pgdb"
	"anti_bruteforce/pkg/redisdb"
)

func main() {
	cfg := config.MustConfig()
	limit := internal.FromConfig(cfg)
	log := logger.GetLogger("anti_bruteforce", cfg.JSONFormat)
	log.Info().Msg(fmt.Sprintf("Config: %+v", cfg))

	pg := pgdb.MustPostgresDB(cfg.PgDSN)
	defer pg.Close()

	redis := redisdb.MustRedisDB(cfg.RedisDSN)
	defer redis.Close()

	storagePG := storage.NewPgStorage(log, pg)
	bucketRedis := bucket.NewLeakyBucket(log, limit, redis)

	uc := usecase.NewAppUseCase(log, storagePG, bucketRedis)

	server := grpc.NewGrpcServer(log, cfg.Port, uc)

	log.Info().Msg(fmt.Sprintf("server listen and serve %s", server.GetAddr()))
	server.Run()

	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT,
	)
	defer cancel()

	select {
	case <-ctx.Done():
		log.Info().Msg("interrupt signal received")
	case err := <-server.Notify():
		log.Error().Err(err).Msg("server error")
	}

	if err := server.Shutdown(); err != nil {
		log.Error().Err(err).Msg("failed to shutdown server")
	}
}
