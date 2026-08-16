package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/dto"
	"github.com/lp/hanfu-rental/internal/model"
	"github.com/lp/hanfu-rental/internal/util"
)

type fakeUserRepo struct {
	users map[string]*model.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: map[string]*model.User{}}
}

func (f *fakeUserRepo) Create(_ context.Context, u *model.User) error {
	f.users[u.Phone] = u
	return nil
}

func (f *fakeUserRepo) FindByPhone(_ context.Context, phone string) (*model.User, error) {
	if u, ok := f.users[phone]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, util.ErrNotFound
}

func (f *fakeUserRepo) FindByID(_ context.Context, id uint) (*model.User, error) {
	for _, u := range f.users {
		if u.ID == id {
			cp := *u
			return &cp, nil
		}
	}
	return nil, util.ErrNotFound
}

func (f *fakeUserRepo) UpdateProfile(_ context.Context, id uint, nickname, avatar string) error {
	u, err := f.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	f.users[u.Phone].Nickname = nickname
	f.users[u.Phone].Avatar = avatar
	return nil
}

func (f *fakeUserRepo) AddDeposit(_ context.Context, id uint, delta float64) error {
	u, err := f.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	f.users[u.Phone].Deposit += delta
	return nil
}

func (f *fakeUserRepo) Count(context.Context) (int64, error) { return int64(len(f.users)), nil }

func newTestUserService(repo UserRepository) *UserService {
	return NewUserService(repo, "test-secret", 72, slog.Default())
}

func TestUserServiceRegister(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newTestUserService(repo)
	tests := []struct {
		name    string
		phone   string
		wantErr bool
		errCode int
	}{
		{name: "valid phone", phone: "13800001111", wantErr: false},
		{name: "duplicate phone", phone: "13800001111", wantErr: true, errCode: constants.CodeConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &dto.RegisterRequest{Phone: tt.phone, Password: "123456", Nickname: "测试顾客"}
			_, err := svc.Register(context.Background(), req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				var appErr *util.AppError
				if errors.As(err, &appErr) && appErr.Code != tt.errCode {
					t.Fatalf("expected code %d, got %d", tt.errCode, appErr.Code)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestUserServiceLogin(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newTestUserService(repo)
	_, _ = svc.Register(context.Background(), &dto.RegisterRequest{Phone: "13800002222", Password: "123456", Nickname: "登录测试"})
	tests := []struct {
		name     string
		phone    string
		password string
		wantErr  bool
	}{
		{name: "correct password", phone: "13800002222", password: "123456", wantErr: false},
		{name: "wrong password", phone: "13800002222", password: "wrong", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.Login(context.Background(), &dto.LoginRequest{Phone: tt.phone, Password: tt.password})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.Token == "" || resp.User == nil {
				t.Fatalf("expected token and user in response")
			}
		})
	}
}

func TestUserServicePayDeposit(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newTestUserService(repo)
	u, _ := svc.Register(context.Background(), &dto.RegisterRequest{Phone: "13800003333", Password: "123456", Nickname: "押金测试"})
	_, err := svc.PayDeposit(context.Background(), u.ID, 300)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	updated, _ := svc.GetProfile(context.Background(), u.ID)
	if updated.Deposit != 300 {
		t.Fatalf("expected deposit 300, got %f", updated.Deposit)
	}
}
