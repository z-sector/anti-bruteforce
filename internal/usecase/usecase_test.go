package usecase

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"anti_bruteforce/internal"
	"anti_bruteforce/internal/delivery/grpc/pb"
	mock_usecase "anti_bruteforce/internal/usecase/mocks"
)

func TestAppUseCase_AddToBlackList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStorage, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		ipNet := getIPNet(t)
		mockStorage.EXPECT().
			ExistsSubnetInBlackList(gomock.Any(), gomock.AssignableToTypeOf(&net.IPNet{})).
			Return(false, nil).
			Times(1)
		mockStorage.EXPECT().
			CreateSubnetInBlackList(gomock.Any(), gomock.AssignableToTypeOf(&net.IPNet{})).
			Return(nil).
			Times(1)

		err := uc.AddToBlackList(context.Background(), ipNet.String())

		require.NoError(t, err)
	})

	t.Run("invalid arguments", func(t *testing.T) {
		_, _, uc := getMockStorageAndMockBucketAndUseCase(t)

		err := uc.AddToBlackList(context.Background(), "123")

		require.ErrorIs(t, err, internal.ErrInvalidArgs)
	})

	t.Run("subnet already exists", func(t *testing.T) {
		mockStorage, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		ipNet := getIPNet(t)
		mockStorage.EXPECT().
			ExistsSubnetInBlackList(gomock.Any(), gomock.AssignableToTypeOf(&net.IPNet{})).
			Return(true, nil).
			Times(1)

		err := uc.AddToBlackList(context.Background(), ipNet.String())

		require.ErrorIs(t, err, internal.ErrBlackListExists)
	})
}

func TestAppUseCase_RemoveFromBlackList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStorage, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		ipNet := getIPNet(t)
		mockStorage.EXPECT().
			ExistsSubnetInBlackList(gomock.Any(), gomock.AssignableToTypeOf(&net.IPNet{})).
			Return(true, nil).
			Times(1)
		mockStorage.EXPECT().
			DeleteSubnetInBlackList(gomock.Any(), gomock.AssignableToTypeOf(&net.IPNet{})).
			Return(nil).
			Times(1)

		err := uc.RemoveFromBlackList(context.Background(), ipNet.String())

		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockStorage, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		ipNet := getIPNet(t)
		mockStorage.EXPECT().
			ExistsSubnetInBlackList(gomock.Any(), gomock.AssignableToTypeOf(&net.IPNet{})).
			Return(false, nil).
			Times(1)

		err := uc.RemoveFromBlackList(context.Background(), ipNet.String())

		require.ErrorIs(t, err, internal.ErrBlackListNotFound)
	})
}

func TestAppUseCase_AddToWhiteList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStorage, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		ipNet := getIPNet(t)
		mockStorage.EXPECT().
			ExistsSubnetInWhiteList(gomock.Any(), gomock.AssignableToTypeOf(&net.IPNet{})).
			Return(false, nil).
			Times(1)
		mockStorage.EXPECT().
			CreateSubnetInWhiteList(gomock.Any(), gomock.AssignableToTypeOf(&net.IPNet{})).
			Return(nil).
			Times(1)

		err := uc.AddToWhiteList(context.Background(), ipNet.String())

		require.NoError(t, err)
	})

	t.Run("invalid arguments", func(t *testing.T) {
		_, _, uc := getMockStorageAndMockBucketAndUseCase(t)

		err := uc.AddToWhiteList(context.Background(), "123")

		require.ErrorIs(t, err, internal.ErrInvalidArgs)
	})

	t.Run("subnet already exists", func(t *testing.T) {
		mockStorage, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		ipNet := getIPNet(t)
		mockStorage.EXPECT().
			ExistsSubnetInWhiteList(gomock.Any(), gomock.AssignableToTypeOf(&net.IPNet{})).
			Return(true, nil).
			Times(1)

		err := uc.AddToWhiteList(context.Background(), ipNet.String())

		require.ErrorIs(t, err, internal.ErrWhiteListExists)
	})
}

func TestAppUseCase_RemoveFromWhiteList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStorage, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		ipNet := getIPNet(t)
		mockStorage.EXPECT().
			ExistsSubnetInWhiteList(gomock.Any(), gomock.AssignableToTypeOf(&net.IPNet{})).
			Return(true, nil).
			Times(1)
		mockStorage.EXPECT().
			DeleteSubnetInWhiteList(gomock.Any(), gomock.AssignableToTypeOf(&net.IPNet{})).
			Return(nil).
			Times(1)

		err := uc.RemoveFromWhiteList(context.Background(), ipNet.String())

		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockStorage, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		ipNet := getIPNet(t)
		mockStorage.EXPECT().
			ExistsSubnetInWhiteList(gomock.Any(), gomock.AssignableToTypeOf(&net.IPNet{})).
			Return(false, nil).
			Times(1)

		err := uc.RemoveFromWhiteList(context.Background(), ipNet.String())

		require.ErrorIs(t, err, internal.ErrWhiteListNotFound)
	})
}

func TestAppUseCase_ClearLists(t *testing.T) {
	expectedErr := errors.New("error")
	mockStorage, _, uc := getMockStorageAndMockBucketAndUseCase(t)
	mockStorage.EXPECT().ClearLists(gomock.Any()).Return(expectedErr).Times(1)

	err := uc.ClearLists(context.Background())

	require.ErrorIs(t, err, expectedErr)
}

func TestAppUseCase_Reset(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, mockBucket, uc := getMockStorageAndMockBucketAndUseCase(t)
		request := getRequest(t)
		mockBucket.EXPECT().ResetLogin(gomock.Any(), request.GetLogin()).Return(nil).Times(1)
		mockBucket.EXPECT().ResetPassword(gomock.Any(), request.GetPassword()).Return(nil).Times(1)
		mockBucket.EXPECT().ResetIP(gomock.Any(), net.ParseIP(request.Ip)).Return(nil).Times(1)

		err := uc.Reset(context.Background(), request)

		require.NoError(t, err)
	})

	t.Run("invalid ip", func(t *testing.T) {
		_, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		request := getRequest(t)
		request.Ip = "111"

		err := uc.Reset(context.Background(), request)

		require.ErrorIs(t, err, internal.ErrInvalidIP)
	})
}

func TestAppUseCase_CheckAuth(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStorage, mockBucket, uc := getMockStorageAndMockBucketAndUseCase(t)
		request := getRequest(t)
		gomock.InOrder(
			mockStorage.EXPECT().
				CheckInWhiteList(gomock.Any(), gomock.AssignableToTypeOf(net.IP{})).
				Return(false, nil).
				Times(1),
			mockStorage.EXPECT().
				CheckInBlackList(gomock.Any(), gomock.AssignableToTypeOf(net.IP{})).
				Return(false, nil).
				Times(1),
		)
		mockBucket.EXPECT().CheckLogin(gomock.Any(), request.GetLogin()).Return(true, nil).Times(1)
		mockBucket.EXPECT().CheckPassword(gomock.Any(), request.GetPassword()).Return(true, nil).Times(1)
		mockBucket.EXPECT().CheckIP(gomock.Any(), gomock.AssignableToTypeOf(net.IP{})).Return(true, nil).Times(1)

		res, err := uc.CheckAuth(context.Background(), request)

		require.NoError(t, err)
		require.True(t, res)
	})

	t.Run("required fields error", func(t *testing.T) {
		_, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		request := getRequest(t)
		request.Login = ""

		_, err := uc.CheckAuth(context.Background(), request)

		require.ErrorIs(t, err, internal.ErrInvalidArgs)
	})

	t.Run("invalid ip", func(t *testing.T) {
		_, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		request := getRequest(t)
		request.Ip = "123"

		_, err := uc.CheckAuth(context.Background(), request)

		require.ErrorIs(t, err, internal.ErrInvalidIP)
	})

	t.Run("is in whitelist", func(t *testing.T) {
		mockStorage, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		request := getRequest(t)
		mockStorage.EXPECT().
			CheckInWhiteList(gomock.Any(), gomock.AssignableToTypeOf(net.IP{})).
			Return(true, nil).
			Times(1)

		res, err := uc.CheckAuth(context.Background(), request)

		require.NoError(t, err)
		require.True(t, res)
	})

	t.Run("is in blacklist", func(t *testing.T) {
		mockStorage, _, uc := getMockStorageAndMockBucketAndUseCase(t)
		request := getRequest(t)
		gomock.InOrder(
			mockStorage.EXPECT().
				CheckInWhiteList(gomock.Any(), gomock.AssignableToTypeOf(net.IP{})).
				Return(false, nil).
				Times(1),
			mockStorage.EXPECT().
				CheckInBlackList(gomock.Any(), gomock.AssignableToTypeOf(net.IP{})).
				Return(true, nil).
				Times(1),
		)

		res, err := uc.CheckAuth(context.Background(), request)

		require.NoError(t, err)
		require.False(t, res)
	})

	t.Run("is in bucket by password", func(t *testing.T) {
		mockStorage, mockBucket, uc := getMockStorageAndMockBucketAndUseCase(t)
		request := getRequest(t)
		gomock.InOrder(
			mockStorage.EXPECT().
				CheckInWhiteList(gomock.Any(), gomock.AssignableToTypeOf(net.IP{})).
				Return(false, nil).
				Times(1),
			mockStorage.EXPECT().
				CheckInBlackList(gomock.Any(), gomock.AssignableToTypeOf(net.IP{})).
				Return(false, nil).
				Times(1),
		)
		mockBucket.EXPECT().CheckLogin(gomock.Any(), request.GetLogin()).Return(true, nil).Times(1)
		mockBucket.EXPECT().CheckPassword(gomock.Any(), request.GetPassword()).Return(false, nil).Times(1)
		mockBucket.EXPECT().CheckIP(gomock.Any(), gomock.AssignableToTypeOf(net.IP{})).Return(true, nil).Times(1)

		res, err := uc.CheckAuth(context.Background(), request)

		require.NoError(t, err)
		require.False(t, res)
	})
}

func getIPNet(t *testing.T) *net.IPNet {
	t.Helper()

	_, ipNet, err := net.ParseCIDR("192.168.0.1/32")
	require.NoError(t, err)
	return ipNet
}

func getRequest(t *testing.T) *pb.AuthCheckRequest {
	t.Helper()

	ipNet := getIPNet(t)
	return &pb.AuthCheckRequest{
		Login:    "login",
		Password: "password",
		Ip:       ipNet.IP.String(),
	}
}

func getMockStorageAndMockBucketAndUseCase(t *testing.T) (
	*mock_usecase.MockStorageI,
	*mock_usecase.MockLeakyBucketI,
	*AppUseCase,
) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockStorage := mock_usecase.NewMockStorageI(ctrl)
	mockBucket := mock_usecase.NewMockLeakyBucketI(ctrl)
	uc := NewAppUseCase(zerolog.Nop(), mockStorage, mockBucket)
	return mockStorage, mockBucket, uc
}
