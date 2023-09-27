package bucket

import (
	"context"
	"fmt"
	"net"

	"github.com/go-redis/redis_rate/v10"
	"github.com/rs/zerolog"

	"anti_bruteforce/internal"
	"anti_bruteforce/pkg/redisdb"
)

const (
	prefKeyLogin    = "login"
	prefKeyPassword = "password"
	prefKeyIP       = "ip"
)

type LeakyBucket struct {
	log           zerolog.Logger
	limiter       *redis_rate.Limiter
	loginLimit    redis_rate.Limit
	passwordLimit redis_rate.Limit
	ipLimit       redis_rate.Limit
}

func NewLeakyBucket(log zerolog.Logger, limit *internal.LimitSettings, redis *redisdb.RedisDB) *LeakyBucket {
	return &LeakyBucket{
		log:     log,
		limiter: redis_rate.NewLimiter(redis.Client),
		loginLimit: redis_rate.Limit{
			Rate:   limit.Login.Burst,
			Burst:  limit.Login.Burst,
			Period: limit.Login.Period,
		},
		passwordLimit: redis_rate.Limit{
			Rate:   limit.Password.Burst,
			Burst:  limit.Password.Burst,
			Period: limit.Password.Period,
		},
		ipLimit: redis_rate.Limit{
			Rate:   limit.IP.Burst,
			Burst:  limit.IP.Burst,
			Period: limit.IP.Period,
		},
	}
}

func (b *LeakyBucket) CheckLogin(ctx context.Context, login string) (bool, error) {
	return b.allow(ctx, b.fmtLoginKey(login), b.loginLimit)
}

func (b *LeakyBucket) ResetLogin(ctx context.Context, login string) error {
	return b.reset(ctx, b.fmtLoginKey(login))
}

func (b *LeakyBucket) CheckPassword(ctx context.Context, pwd string) (bool, error) {
	return b.allow(ctx, b.fmtPasswordKey(pwd), b.passwordLimit)
}

func (b *LeakyBucket) ResetPassword(ctx context.Context, pwd string) error {
	return b.reset(ctx, b.fmtPasswordKey(pwd))
}

func (b *LeakyBucket) CheckIP(ctx context.Context, ip net.IP) (bool, error) {
	return b.allow(ctx, b.fmtIPKey(ip), b.ipLimit)
}

func (b *LeakyBucket) ResetIP(ctx context.Context, ip net.IP) error {
	return b.reset(ctx, b.fmtIPKey(ip))
}

func (b *LeakyBucket) fmtLoginKey(login string) string {
	return fmt.Sprintf("%s:%s", prefKeyLogin, login)
}

func (b *LeakyBucket) fmtPasswordKey(pwd string) string {
	return fmt.Sprintf("%s:%s", prefKeyPassword, pwd)
}

func (b *LeakyBucket) fmtIPKey(ip net.IP) string {
	return fmt.Sprintf("%s:%s", prefKeyIP, ip.String())
}

func (b *LeakyBucket) allow(ctx context.Context, key string, limit redis_rate.Limit) (bool, error) {
	res, err := b.limiter.Allow(ctx, key, limit)
	if err != nil {
		return false, fmt.Errorf("LeakyBucket - allow - limiter.Allow: %w", err)
	}
	b.log.Info().
		Int("allowed", res.Allowed).
		Int("remaining", res.Remaining).
		Dur("resetAfter", res.ResetAfter).
		Dur("retryAfter", res.RetryAfter).
		Msg(key)
	if res.Allowed == 0 {
		return false, nil
	}
	return true, nil
}

func (b *LeakyBucket) reset(ctx context.Context, key string) error {
	if err := b.limiter.Reset(ctx, key); err != nil {
		return fmt.Errorf("LeakyBucket - reset - limiter.Reset: %w", err)
	}
	return nil
}
