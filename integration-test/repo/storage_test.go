//go:build integration

package repo_test

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"anti_bruteforce/internal/repositories/storage"
	"anti_bruteforce/pkg/pgdb"
)

func TestStorage(t *testing.T) {
	suite.Run(t, new(storageTestSuite))
}

type storageTestSuite struct {
	suite.Suite
	pgClient       *pgdb.PostgresDB
	Storage        *storage.PgStorage
	random         *rand.Rand
	ctx            context.Context
	existsBlackSQL string
	insertBlackSQL string
	existsWhiteSQL string
	insertWhiteSQL string
}

func (s *storageTestSuite) SetupSuite() {
	dsn := os.Getenv("APP_PG_DSN")
	s.Require().NotEmpty(dsn)
	s.pgClient = pgdb.MustPostgresDB(dsn)
	s.Storage = storage.NewPgStorage(zerolog.Nop(), s.pgClient)
	s.ctx = context.Background()
	s.existsBlackSQL = `SELECT EXISTS (SELECT 1 FROM black_list WHERE ip_address = $1)`
	s.existsWhiteSQL = `SELECT EXISTS (SELECT 1 FROM white_list WHERE ip_address = $1)`
	s.insertBlackSQL = `INSERT INTO black_list (ip_address) VALUES ($1)`
	s.insertWhiteSQL = `INSERT INTO white_list (ip_address) VALUES ($1)`

	seed := time.Now().UnixNano()
	s.T().Logf("rand seed: %d\n", seed)
	s.random = rand.New(rand.NewSource(seed))
}

func (s *storageTestSuite) TearDownSuite() {
	s.cleanUp()
	s.pgClient.Close()
}

func (s *storageTestSuite) TearDownTest() {
	s.cleanUp()
}

func (s *storageTestSuite) TearDownSubTest() {
	s.cleanUp()
}

func (s *storageTestSuite) cleanUp() {
	_, err := s.pgClient.Pool.Exec(s.ctx, `TRUNCATE TABLE black_list, white_list RESTART IDENTITY`)
	s.Require().NoError(err)
}

func (s *storageTestSuite) TestCreateSubnetInBlackList() {
	ipNet := randomSimpleIPNet(s.T(), s.random)

	err := s.Storage.CreateSubnetInBlackList(s.ctx, ipNet)

	s.Require().NoError(err)
	var found bool
	err = s.pgClient.Pool.QueryRow(s.ctx, s.existsBlackSQL, ipNet.String()).Scan(&found)
	s.Require().NoError(err)
	s.Require().True(found)
}

func (s *storageTestSuite) TestExistsSubnetInBlackList() {
	testcases := []struct {
		fixture *net.IPNet
		search  *net.IPNet
		wont    bool
	}{
		{
			fixture: &net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.CIDRMask(32, 32),
			},
			search: &net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.CIDRMask(32, 32),
			},
			wont: true,
		},
		{
			fixture: &net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.CIDRMask(32, 32),
			},
			search: &net.IPNet{
				IP:   net.IP{192, 168, 1, 1},
				Mask: net.CIDRMask(32, 32),
			},
			wont: false,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		s.Run(fmt.Sprintf("case №%d:%t", i, tc.wont), func() {
			_, err := s.pgClient.Pool.Exec(s.ctx, s.insertBlackSQL, tc.fixture.String())
			s.Require().NoError(err)

			res, err := s.Storage.ExistsSubnetInBlackList(s.ctx, tc.search)

			s.Require().NoError(err)
			s.Require().Equal(tc.wont, res)
		})
	}
}

func (s *storageTestSuite) TestDeleteSubnetInBlackList() {
	ipNet := randomSimpleIPNet(s.T(), s.random)
	_, err := s.pgClient.Pool.Exec(s.ctx, s.insertBlackSQL, ipNet.String())
	s.Require().NoError(err)

	err = s.Storage.DeleteSubnetInBlackList(s.ctx, ipNet)

	s.Require().NoError(err)
	var found bool
	err = s.pgClient.Pool.QueryRow(s.ctx, s.existsBlackSQL, ipNet.String()).Scan(&found)
	s.Require().NoError(err)
	s.Require().False(found)
}

func (s *storageTestSuite) TestCheckInBlackList() {
	testcases := []struct {
		ipNet *net.IPNet
		ip    net.IP
		wont  bool
	}{
		{
			ipNet: &net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.CIDRMask(30, 32),
			},
			ip:   net.IP{192, 168, 1, 1},
			wont: true,
		},
		{
			ipNet: &net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.CIDRMask(30, 32),
			},
			ip:   net.IP{192, 168, 1, 10},
			wont: false,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		s.Run(fmt.Sprintf("case №%d:%t", i, tc.wont), func() {
			_, err := s.pgClient.Pool.Exec(s.ctx, s.insertBlackSQL, tc.ipNet.String())
			s.Require().NoError(err)

			res, err := s.Storage.CheckInBlackList(s.ctx, tc.ip)

			s.Require().NoError(err)
			s.Require().Equal(tc.wont, res)
		})
	}
}

func (s *storageTestSuite) TestClearLists() {
	ipNetWhite := randomSimpleIPNet(s.T(), s.random)
	ipNetBlack := randomSimpleIPNet(s.T(), s.random)
	_, err := s.pgClient.Pool.Exec(s.ctx, s.insertWhiteSQL, ipNetWhite.String())
	s.Require().NoError(err)
	_, err = s.pgClient.Pool.Exec(s.ctx, s.insertBlackSQL, ipNetBlack.String())
	s.Require().NoError(err)

	err = s.Storage.ClearLists(s.ctx)

	s.Require().NoError(err)
	var found bool
	err = s.pgClient.Pool.QueryRow(s.ctx, `SELECT EXISTS (SELECT 1 FROM white_list)`).Scan(&found)
	s.Require().NoError(err)
	s.Require().False(found)
	err = s.pgClient.Pool.QueryRow(s.ctx, `SELECT EXISTS (SELECT 1 FROM black_list)`).Scan(&found)
	s.Require().NoError(err)
	s.Require().False(found)
}

func (s *storageTestSuite) TestCreateSubnetInWhiteList() {
	ipNet := randomSimpleIPNet(s.T(), s.random)

	err := s.Storage.CreateSubnetInWhiteList(s.ctx, ipNet)

	s.Require().NoError(err)
	var found bool
	err = s.pgClient.Pool.QueryRow(s.ctx, s.existsWhiteSQL, ipNet.String()).Scan(&found)
	s.Require().NoError(err)
	s.Require().True(found)
}

func (s *storageTestSuite) TestExistsSubnetInWhiteList() {
	testcases := []struct {
		fixture *net.IPNet
		search  *net.IPNet
		wont    bool
	}{
		{
			fixture: &net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.CIDRMask(32, 32),
			},
			search: &net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.CIDRMask(32, 32),
			},
			wont: true,
		},
		{
			fixture: &net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.CIDRMask(32, 32),
			},
			search: &net.IPNet{
				IP:   net.IP{192, 168, 1, 1},
				Mask: net.CIDRMask(32, 32),
			},
			wont: false,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		s.Run(fmt.Sprintf("case №%d:%t", i, tc.wont), func() {
			_, err := s.pgClient.Pool.Exec(s.ctx, s.insertWhiteSQL, tc.fixture.String())
			s.Require().NoError(err)

			res, err := s.Storage.ExistsSubnetInWhiteList(s.ctx, tc.search)

			s.Require().NoError(err)
			s.Require().Equal(tc.wont, res)
		})
	}
}

func (s *storageTestSuite) TestDeleteSubnetInWhiteList() {
	ipNet := randomSimpleIPNet(s.T(), s.random)
	_, err := s.pgClient.Pool.Exec(s.ctx, s.insertWhiteSQL, ipNet.String())
	s.Require().NoError(err)

	err = s.Storage.DeleteSubnetInWhiteList(s.ctx, ipNet)

	s.Require().NoError(err)
	var found bool
	err = s.pgClient.Pool.QueryRow(s.ctx, s.existsWhiteSQL, ipNet.String()).Scan(&found)
	s.Require().NoError(err)
	s.Require().False(found)
}

func (s *storageTestSuite) TestCheckInWhiteList() {
	testcases := []struct {
		ipNet *net.IPNet
		ip    net.IP
		wont  bool
	}{
		{
			ipNet: &net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.CIDRMask(30, 32),
			},
			ip:   net.IP{192, 168, 1, 1},
			wont: true,
		},
		{
			ipNet: &net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.CIDRMask(30, 32),
			},
			ip:   net.IP{192, 168, 1, 10},
			wont: false,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		s.Run(fmt.Sprintf("case №%d:%t", i, tc.wont), func() {
			_, err := s.pgClient.Pool.Exec(s.ctx, s.insertWhiteSQL, tc.ipNet.String())
			s.Require().NoError(err)

			res, err := s.Storage.CheckInWhiteList(s.ctx, tc.ip)

			s.Require().NoError(err)
			s.Require().Equal(tc.wont, res)
		})
	}
}

func randomSimpleIPNet(t *testing.T, random *rand.Rand) *net.IPNet {
	t.Helper()

	buf := make([]byte, 4)
	n, err := random.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 4, n)
	return &net.IPNet{
		IP:   net.IPv4(buf[0], buf[1], buf[2], buf[3]),
		Mask: net.CIDRMask(32, 32),
	}
}
