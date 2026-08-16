package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/dto"
	"github.com/lp/hanfu-rental/internal/model"
	"github.com/lp/hanfu-rental/internal/util"
)

type nilUserRepo struct{}

func (m *nilUserRepo) Create(ctx context.Context, u *model.User) error { return nil }
func (m *nilUserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) { return nil, nil }
func (m *nilUserRepo) FindByID(ctx context.Context, id uint) (*model.User, error) { return nil, nil }
func (m *nilUserRepo) UpdateProfile(ctx context.Context, id uint, nickname, avatar string) error { return nil }
func (m *nilUserRepo) AddDeposit(ctx context.Context, id uint, delta float64) error { return nil }
func (m *nilUserRepo) Count(ctx context.Context) (int64, error) { return 0, nil }

func newNilUserService() *UserService {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	return NewUserService(&nilUserRepo{}, "test-secret", 72, logger)
}

func TestGetProfileNilReturnsError(t *testing.T) {
	svc := newNilUserService()
	user, err := svc.GetProfile(context.Background(), 1)
	if err == nil {
		t.Fatalf("expected error for nil user, got user=%+v", user)
	}
}

func TestLoginNilReturnsUnauthorized(t *testing.T) {
	svc := newNilUserService()
	_, err := svc.Login(context.Background(), &dto.LoginRequest{Phone: "13800009999", Password: "123456"})
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.Code != constants.CodeUnauthorized {
		t.Fatalf("expected CodeUnauthorized, got %v", err)
	}
}

func TestDepositFlow(t *testing.T) {
	svc := newTestUserService(newFakeUserRepo())
	u, _ := svc.Register(context.Background(), &dto.RegisterRequest{Phone: "13800001111", Password: "123456", Nickname: "测试"})
	u, _ = svc.PayDeposit(context.Background(), u.ID, 300)
	if u.Deposit != 300 {
		t.Fatalf("deposit after pay = %f, want 300", u.Deposit)
	}
	u, _ = svc.RefundDeposit(context.Background(), u.ID, 100)
	if u.Deposit != 200 {
		t.Fatalf("deposit after refund = %f, want 200", u.Deposit)
	}
}
