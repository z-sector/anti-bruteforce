//go:build integration

package integration_test

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	"anti_bruteforce/internal"
	"anti_bruteforce/internal/repositories/bucket"
	"anti_bruteforce/pkg/redisdb"
)

func TestBucket(t *testing.T) {
	suite.Run(t, new(bucketTestSuite))
}

type bucketTestSuite struct {
	suite.Suite
	redis       *redisdb.RedisDB
	limit       *internal.LimitSettings
	LeakyBucket *bucket.LeakyBucket
	ctx         context.Context
}

func (s *bucketTestSuite) SetupSuite() {
	dsn := os.Getenv("APP_REDIS_DSN")
	s.Require().NotEmpty(dsn)
	s.redis = redisdb.MustRedisDB(dsn)
	limitItem := internal.LimitItem{
		Burst:  1,
		Period: 100 * time.Millisecond,
	}
	s.limit = &internal.LimitSettings{
		Login:    limitItem,
		IP:       limitItem,
		Password: limitItem,
	}
	s.LeakyBucket = bucket.NewLeakyBucket(zerolog.Nop(), s.limit, s.redis)
	s.ctx = context.Background()
}

func (s *bucketTestSuite) TearDownSuite() {
	s.cleanUp()
	s.redis.Close()
}

func (s *bucketTestSuite) TearDownTest() {
	s.cleanUp()
}

func (s *bucketTestSuite) cleanUp() {
	err := s.redis.Client.FlushDB(s.ctx).Err()
	s.Require().NoError(err)
}

func (s *bucketTestSuite) TestLimitLogin() {
	login := "login"

	ok, err := s.LeakyBucket.CheckLogin(s.ctx, login)

	s.Require().NoError(err)
	s.Require().True(ok)

	ok, err = s.LeakyBucket.CheckLogin(s.ctx, login)

	s.Require().NoError(err)
	s.Require().False(ok)

	err = s.LeakyBucket.ResetLogin(s.ctx, login)

	s.Require().NoError(err)

	ok, err = s.LeakyBucket.CheckLogin(s.ctx, login)

	s.Require().NoError(err)
	s.Require().True(ok)

	time.Sleep(s.limit.Login.Period)

	ok, err = s.LeakyBucket.CheckLogin(s.ctx, login)

	s.Require().NoError(err)
	s.Require().True(ok)
}

func (s *bucketTestSuite) TestLimitPassword() {
	password := "password"

	ok, err := s.LeakyBucket.CheckPassword(s.ctx, password)

	s.Require().NoError(err)
	s.Require().True(ok)

	ok, err = s.LeakyBucket.CheckPassword(s.ctx, password)

	s.Require().NoError(err)
	s.Require().False(ok)

	err = s.LeakyBucket.ResetPassword(s.ctx, password)

	s.Require().NoError(err)

	ok, err = s.LeakyBucket.CheckPassword(s.ctx, password)

	s.Require().NoError(err)
	s.Require().True(ok)

	time.Sleep(s.limit.Password.Period)

	ok, err = s.LeakyBucket.CheckPassword(s.ctx, password)

	s.Require().NoError(err)
	s.Require().True(ok)
}

func (s *bucketTestSuite) TestLimitIP() {
	ip := net.IP{192, 168, 1, 0}

	ok, err := s.LeakyBucket.CheckIP(s.ctx, ip)

	s.Require().NoError(err)
	s.Require().True(ok)

	ok, err = s.LeakyBucket.CheckIP(s.ctx, ip)

	s.Require().NoError(err)
	s.Require().False(ok)

	err = s.LeakyBucket.ResetIP(s.ctx, ip)

	s.Require().NoError(err)

	ok, err = s.LeakyBucket.CheckIP(s.ctx, ip)

	s.Require().NoError(err)
	s.Require().True(ok)

	time.Sleep(s.limit.IP.Period)

	ok, err = s.LeakyBucket.CheckIP(s.ctx, ip)

	s.Require().NoError(err)
	s.Require().True(ok)
}
