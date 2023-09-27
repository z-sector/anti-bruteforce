package grpc

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"

	mock_grpc "anti_bruteforce/internal/delivery/grpc/mocks"
	"anti_bruteforce/internal/delivery/grpc/pb"
	"anti_bruteforce/internal/models"
)

func TestHandlerGrpc_AddToBlackList(t *testing.T) {
	testcases := []struct {
		expectedErr error
	}{
		{
			expectedErr: errors.New("error"),
		},
		{
			expectedErr: nil,
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(fmt.Sprintf("case №%d", i), func(t *testing.T) {
			mockUC, hdl := getMockUCAndHandler(t)
			in := getSubnetAddress(t)
			mockUC.EXPECT().AddToBlackList(gomock.Any(), in.GetSubnetAddress()).Return(tc.expectedErr).Times(1)

			res, err := hdl.AddToBlackList(context.Background(), in)

			require.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				require.NotNil(t, res)
			}
		})
	}
}

func TestHandlerGrpc_RemoveFromBlackList(t *testing.T) {
	testcases := []struct {
		expectedErr error
	}{
		{
			expectedErr: errors.New("error"),
		},
		{
			expectedErr: nil,
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(fmt.Sprintf("case №%d", i), func(t *testing.T) {
			mockUC, hdl := getMockUCAndHandler(t)
			in := getSubnetAddress(t)
			mockUC.EXPECT().RemoveFromBlackList(gomock.Any(), in.GetSubnetAddress()).Return(tc.expectedErr).Times(1)

			res, err := hdl.RemoveFromBlackList(context.Background(), in)

			require.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				require.NotNil(t, res)
			}
		})
	}
}

func TestHandlerGrpc_AddToWhiteList(t *testing.T) {
	testcases := []struct {
		expectedErr error
	}{
		{
			expectedErr: errors.New("error"),
		},
		{
			expectedErr: nil,
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(fmt.Sprintf("case №%d", i), func(t *testing.T) {
			mockUC, hdl := getMockUCAndHandler(t)
			in := getSubnetAddress(t)
			mockUC.EXPECT().AddToWhiteList(gomock.Any(), in.GetSubnetAddress()).Return(tc.expectedErr).Times(1)

			res, err := hdl.AddToWhiteList(context.Background(), in)

			require.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				require.NotNil(t, res)
			}
		})
	}
}

func TestHandlerGrpc_RemoveFromWhiteList(t *testing.T) {
	testcases := []struct {
		expectedErr error
	}{
		{
			expectedErr: errors.New("error"),
		},
		{
			expectedErr: nil,
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(fmt.Sprintf("case №%d", i), func(t *testing.T) {
			mockUC, hdl := getMockUCAndHandler(t)
			in := getSubnetAddress(t)
			mockUC.EXPECT().RemoveFromWhiteList(gomock.Any(), in.GetSubnetAddress()).Return(tc.expectedErr).Times(1)

			res, err := hdl.RemoveFromWhiteList(context.Background(), in)

			require.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				require.NotNil(t, res)
			}
		})
	}
}

func TestHandlerGrpc_ClearLists(t *testing.T) {
	testcases := []struct {
		expectedErr error
	}{
		{
			expectedErr: errors.New("error"),
		},
		{
			expectedErr: nil,
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(fmt.Sprintf("case №%d", i), func(t *testing.T) {
			mockUC, hdl := getMockUCAndHandler(t)
			mockUC.EXPECT().ClearLists(gomock.Any()).Return(tc.expectedErr).Times(1)

			res, err := hdl.ClearLists(context.Background(), &emptypb.Empty{})

			require.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				require.NotNil(t, res)
			}
		})
	}
}

func TestHandlerGrpc_AuthCheck(t *testing.T) {
	testcases := []struct {
		wontAccepted bool
		expectedErr  error
	}{
		{
			wontAccepted: false,
			expectedErr:  errors.New("error"),
		},
		{
			wontAccepted: false,
			expectedErr:  nil,
		},
		{
			wontAccepted: true,
			expectedErr:  nil,
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(fmt.Sprintf("case №%d", i), func(t *testing.T) {
			mockUC, hdl := getMockUCAndHandler(t)
			request := getAuthCheckRequest(t)
			authCheck := models.AuthCheck{
				Login:    request.GetLogin(),
				Password: request.GetPassword(),
				IP:       request.GetIp(),
			}
			mockUC.EXPECT().CheckAuth(gomock.Any(), authCheck).Return(tc.wontAccepted, tc.expectedErr).Times(1)

			res, err := hdl.AuthCheck(context.Background(), request)

			require.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				require.NotNil(t, res)
				require.Equal(t, tc.wontAccepted, res.Accepted)
			}
		})
	}
}

func TestHandlerGrpc_Reset(t *testing.T) {
	testcases := []struct {
		expectedErr error
	}{
		{
			expectedErr: errors.New("error"),
		},
		{
			expectedErr: nil,
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(fmt.Sprintf("case №%d", i), func(t *testing.T) {
			mockUC, hdl := getMockUCAndHandler(t)
			request := getResetBucketRequest(t)
			data := models.ResetBucketData{
				Login: request.GetLogin(),
				IP:    request.GetIp(),
			}
			mockUC.EXPECT().ResetBucket(gomock.Any(), data).Return(tc.expectedErr).Times(1)

			res, err := hdl.ResetBucket(context.Background(), request)

			require.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				require.NotNil(t, res)
			}
		})
	}
}

func getMockUCAndHandler(t *testing.T) (*mock_grpc.MockUseCaseI, *HandlerGrpc) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockUC := mock_grpc.NewMockUseCaseI(ctrl)
	hdl := NewHandlerGrpc(zerolog.Nop(), mockUC)
	return mockUC, hdl
}

func getSubnetAddress(t *testing.T) *pb.SubnetAddress {
	t.Helper()
	return &pb.SubnetAddress{SubnetAddress: "192.168.0.1/32"}
}

func getAuthCheckRequest(t *testing.T) *pb.AuthCheckRequest {
	t.Helper()
	return &pb.AuthCheckRequest{
		Login:    "login",
		Password: "password",
		Ip:       "192.168.0.5/32",
	}
}

func getResetBucketRequest(t *testing.T) *pb.ResetBucketRequest {
	t.Helper()
	return &pb.ResetBucketRequest{
		Login: "login",
		Ip:    "192.168.0.5/32",
	}
}
