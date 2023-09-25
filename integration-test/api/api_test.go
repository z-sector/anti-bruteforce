//go:build integration

package api_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"anti_bruteforce/internal/delivery/grpc/pb"
	"anti_bruteforce/pkg/pgdb"
	"anti_bruteforce/pkg/redisdb"
)

func TestAPI(t *testing.T) {
	suite.Run(t, new(apiTestSuite))
}

type apiTestSuite struct {
	suite.Suite
	redis    *redisdb.RedisDB
	pgClient *pgdb.PostgresDB
	conn     *grpc.ClientConn
	client   pb.AntiBruteForceServiceClient
	ctx      context.Context
}

func (s *apiTestSuite) SetupSuite() {
	var err error

	redisDSN := os.Getenv("APP_REDIS_DSN")
	s.Require().NotEmpty(redisDSN)
	s.redis = redisdb.MustRedisDB(redisDSN)

	pgDSN := os.Getenv("APP_PG_DSN")
	s.Require().NotEmpty(pgDSN)
	s.pgClient = pgdb.MustPostgresDB(pgDSN)

	serviceHost := os.Getenv("APP_SERVICE_HOST")
	s.Require().NotEmpty(serviceHost)
	s.conn, err = grpc.Dial(serviceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.client = pb.NewAntiBruteForceServiceClient(s.conn)

	s.ctx = context.Background()
}

func (s *apiTestSuite) TearDownSuite() {
	s.cleanUp()
	s.pgClient.Close()
	s.redis.Close()
}

func (s *apiTestSuite) TearDownTest() {
	s.cleanUp()
}

func (s *apiTestSuite) TearDownSubTest() {
	s.cleanUp()
}

func (s *apiTestSuite) cleanUp() {
	_, err := s.pgClient.Pool.Exec(s.ctx, `TRUNCATE TABLE black_list, white_list RESTART IDENTITY`)
	s.Require().NoError(err)
	err = s.redis.Client.FlushDB(s.ctx).Err()
	s.Require().NoError(err)
}

func (s *apiTestSuite) TestBlackList() {
	in := &pb.SubnetAddress{SubnetAddress: "192.168.1.0/32"}

	_, err := s.client.AddToBlackList(s.ctx, in)

	s.Require().NoError(err)

	_, err = s.client.AddToBlackList(s.ctx, in)

	s.Require().Error(err)

	_, err = s.client.RemoveFromBlackList(s.ctx, in)

	s.Require().NoError(err)

	_, err = s.client.RemoveFromBlackList(s.ctx, in)

	s.Require().Error(err)
}

func (s *apiTestSuite) TestWhiteList() {
	in := &pb.SubnetAddress{SubnetAddress: "192.168.2.0/32"}

	_, err := s.client.AddToWhiteList(s.ctx, in)

	s.Require().NoError(err)

	_, err = s.client.AddToWhiteList(s.ctx, in)

	s.Require().Error(err)

	_, err = s.client.RemoveFromWhiteList(s.ctx, in)

	s.Require().NoError(err)

	_, err = s.client.RemoveFromWhiteList(s.ctx, in)

	s.Require().Error(err)
}

func (s *apiTestSuite) TestAuthCheckTrueByWhiteList() {
	in := &pb.SubnetAddress{SubnetAddress: "192.168.3.0/30"}

	_, err := s.client.AddToWhiteList(s.ctx, in)

	s.Require().NoError(err)

	for i := 0; i < 10; i++ {
		res, err := s.client.AuthCheck(s.ctx, &pb.AuthCheckRequest{
			Login:    "login",
			Password: "password",
			Ip:       "192.168.3.1",
		})

		s.Require().NoError(err)
		s.Require().True(res.Accepted)
	}
}

func (s *apiTestSuite) TestAuthCheckFalseByBlackList() {
	in := &pb.SubnetAddress{SubnetAddress: "192.168.4.0/30"}

	_, err := s.client.AddToBlackList(s.ctx, in)

	s.Require().NoError(err)

	for i := 0; i < 10; i++ {
		res, err := s.client.AuthCheck(s.ctx, &pb.AuthCheckRequest{
			Login:    "login",
			Password: "password",
			Ip:       "192.168.4.1",
		})

		s.Require().NoError(err)
		s.Require().False(res.Accepted)
	}
}

func (s *apiTestSuite) TestAuthCheckByIPLimit() {
	request := &pb.AuthCheckRequest{
		Login:    "login",
		Password: "password",
		Ip:       "192.168.5.0",
	}

	res, err := s.client.AuthCheck(s.ctx, request)

	s.Require().NoError(err)
	s.Require().True(res.Accepted)

	for i := 0; i < 10; i++ {
		suf := fmt.Sprintf("%d", i)
		request.Login = request.Login + suf
		request.Password = request.Password + suf

		res, err := s.client.AuthCheck(s.ctx, request)

		s.Require().NoError(err)
		s.Require().False(res.Accepted)
	}

	resetRequest := &pb.ResetBucketRequest{
		Login: "",
		Ip:    request.GetIp(),
	}

	_, err = s.client.ResetBucket(s.ctx, resetRequest)

	s.Require().NoError(err)

	request.Login = request.Login + "new"
	request.Password = request.Password + "new"

	res, err = s.client.AuthCheck(s.ctx, request)

	s.Require().NoError(err)
	s.Require().True(res.Accepted)
}

func (s *apiTestSuite) TestAuthCheckByLoginLimit() {
	rawIP := "192.168.6."
	request := &pb.AuthCheckRequest{
		Login:    "login",
		Password: "password",
		Ip:       rawIP + "0",
	}

	res, err := s.client.AuthCheck(s.ctx, request)

	s.Require().NoError(err)
	s.Require().True(res.Accepted)

	for i := 1; i < 11; i++ {
		suf := fmt.Sprintf("%d", i)
		request.Ip = rawIP + suf
		request.Password = request.Password + suf

		res, err := s.client.AuthCheck(s.ctx, request)

		s.Require().NoError(err)
		s.Require().False(res.Accepted)
	}

	resetRequest := &pb.ResetBucketRequest{
		Login: request.GetLogin(),
		Ip:    "",
	}

	_, err = s.client.ResetBucket(s.ctx, resetRequest)

	s.Require().NoError(err)

	request.Ip = rawIP + "255"
	request.Password = request.Password + "new"

	res, err = s.client.AuthCheck(s.ctx, request)

	s.Require().NoError(err)
	s.Require().True(res.Accepted)
}

func (s *apiTestSuite) TestAuthCheckByPasswordLimit() {
	rawIP := "192.168.7."
	request := &pb.AuthCheckRequest{
		Login:    "login",
		Password: "password",
		Ip:       rawIP + "0",
	}

	res, err := s.client.AuthCheck(s.ctx, request)

	s.Require().NoError(err)
	s.Require().True(res.Accepted)

	for i := 1; i < 11; i++ {
		suf := fmt.Sprintf("%d", i)
		request.Ip = rawIP + suf
		request.Login = request.Login + suf

		res, err := s.client.AuthCheck(s.ctx, request)

		s.Require().NoError(err)
		s.Require().False(res.Accepted)
	}
}
