package usecase

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog"

	"anti_bruteforce/internal"
	"anti_bruteforce/internal/models"
)

type StorageI interface {
	CreateSubnetInBlackList(ctx context.Context, ipNet *net.IPNet) error
	ExistsSubnetInBlackList(ctx context.Context, ipNet *net.IPNet) (bool, error)
	DeleteSubnetInBlackList(ctx context.Context, ipNet *net.IPNet) error
	CreateSubnetInWhiteList(ctx context.Context, ipNet *net.IPNet) error
	ExistsSubnetInWhiteList(ctx context.Context, ipNet *net.IPNet) (bool, error)
	DeleteSubnetInWhiteList(ctx context.Context, ipNet *net.IPNet) error
	ClearLists(ctx context.Context) error
	CheckInBlackList(ctx context.Context, ip net.IP) (bool, error)
	CheckInWhiteList(ctx context.Context, ip net.IP) (bool, error)
}

type LeakyBucketI interface {
	CheckLogin(ctx context.Context, login string) (bool, error)
	ResetLogin(ctx context.Context, login string) error
	CheckPassword(ctx context.Context, pwd string) (bool, error)
	ResetPassword(ctx context.Context, pwd string) error
	CheckIP(ctx context.Context, ip net.IP) (bool, error)
	ResetIP(ctx context.Context, ip net.IP) error
}

type AppUseCase struct {
	log     zerolog.Logger
	storage StorageI
	bucket  LeakyBucketI
}

func NewAppUseCase(log zerolog.Logger, storage StorageI, bucket LeakyBucketI) *AppUseCase {
	return &AppUseCase{log: log, storage: storage, bucket: bucket}
}

func (u *AppUseCase) AddToBlackList(ctx context.Context, subnet string) error {
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return fmt.Errorf("%w: %s", internal.ErrInvalidArgs, err.Error())
	}

	exists, err := u.storage.ExistsSubnetInBlackList(ctx, ipNet)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w (address=%s)", internal.ErrBlackListExists, subnet)
	}

	return u.storage.CreateSubnetInBlackList(ctx, ipNet)
}

func (u *AppUseCase) RemoveFromBlackList(ctx context.Context, subnet string) error {
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return fmt.Errorf("%w: %s", internal.ErrInvalidArgs, err.Error())
	}

	exists, err := u.storage.ExistsSubnetInBlackList(ctx, ipNet)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w (address=%s)", internal.ErrBlackListNotFound, subnet)
	}

	return u.storage.DeleteSubnetInBlackList(ctx, ipNet)
}

func (u *AppUseCase) AddToWhiteList(ctx context.Context, subnet string) error {
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return fmt.Errorf("%w: %s", internal.ErrInvalidArgs, err.Error())
	}

	exists, err := u.storage.ExistsSubnetInWhiteList(ctx, ipNet)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w (address=%s)", internal.ErrWhiteListExists, subnet)
	}

	return u.storage.CreateSubnetInWhiteList(ctx, ipNet)
}

func (u *AppUseCase) RemoveFromWhiteList(ctx context.Context, subnet string) error {
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return fmt.Errorf("%w: %s", internal.ErrInvalidArgs, err.Error())
	}

	exists, err := u.storage.ExistsSubnetInWhiteList(ctx, ipNet)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w (address=%s)", internal.ErrWhiteListNotFound, subnet)
	}

	return u.storage.DeleteSubnetInWhiteList(ctx, ipNet)
}

func (u *AppUseCase) ClearLists(ctx context.Context) error {
	return u.storage.ClearLists(ctx)
}

func (u *AppUseCase) CheckAuth(ctx context.Context, data models.AuthCheck) (bool, error) {
	if err := u.validateReqFields(data); err != nil {
		return false, err
	}

	ip := net.ParseIP(data.IP)
	if nil == ip {
		return false, internal.ErrInvalidIP
	}

	ok, end, err := u.checkInList(ctx, ip)
	if err != nil {
		return false, err
	}
	if end != 0 {
		return ok, nil
	}

	return u.checkInBucket(ctx, data.Login, data.Password, ip)
}

func (u *AppUseCase) ResetBucket(ctx context.Context, data models.ResetBucketData) error {
	ip := net.ParseIP(data.IP)
	if data.IP != "" && nil == ip {
		return internal.ErrInvalidIP
	}

	var err error

	if data.Login != "" {
		err = u.bucket.ResetLogin(ctx, data.Login)
		if err != nil {
			return err
		}
	}

	if ip != nil {
		err = u.bucket.ResetIP(ctx, ip)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *AppUseCase) validateReqFields(data models.AuthCheck) error {
	var errFields []string

	if data.Login == "" {
		errFields = append(errFields, "login")
	}
	if data.Password == "" {
		errFields = append(errFields, "password")
	}
	if data.IP == "" {
		errFields = append(errFields, "ip")
	}

	if errFields == nil {
		return nil
	}

	msg := fmt.Sprintf("fields %v must be specified", errFields)
	return fmt.Errorf("%w: %s", internal.ErrInvalidArgs, msg)
}

func (u *AppUseCase) checkInList(ctx context.Context, ip net.IP) (bool, int, error) {
	isWhite, err := u.storage.CheckInWhiteList(ctx, ip)
	if err != nil {
		return false, 0, err
	}
	if isWhite {
		return true, 1, nil
	}

	isBlack, err := u.storage.CheckInBlackList(ctx, ip)
	if err != nil {
		return false, 0, err
	}
	if isBlack {
		return false, 1, nil
	}

	return false, 0, nil
}

func (u *AppUseCase) checkInBucket(ctx context.Context, login string, pwd string, ip net.IP) (bool, error) {
	okLogin, err := u.bucket.CheckLogin(ctx, login)
	if err != nil {
		return false, err
	}
	okPassword, err := u.bucket.CheckPassword(ctx, pwd)
	if err != nil {
		return false, err
	}
	okIP, err := u.bucket.CheckIP(ctx, ip)
	if err != nil {
		return false, err
	}
	return okLogin && okPassword && okIP, nil
}
