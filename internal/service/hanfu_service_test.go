package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/dto"
	"github.com/lp/hanfu-rental/internal/model"
	"github.com/lp/hanfu-rental/internal/util"
)

type fakeHanfuRepo struct {
	hanfus map[uint]*model.Hanfu
	nextID uint
}

func newFakeHanfuRepo() *fakeHanfuRepo {
	return &fakeHanfuRepo{hanfus: map[uint]*model.Hanfu{}, nextID: 1}
}

func (f *fakeHanfuRepo) Create(_ context.Context, h *model.Hanfu) error {
	h.ID = f.nextID
	f.nextID++
	f.hanfus[h.ID] = h
	return nil
}

func (f *fakeHanfuRepo) FindByID(_ context.Context, id uint) (*model.Hanfu, error) {
	if h, ok := f.hanfus[id]; ok {
		cp := *h
		return &cp, nil
	}
	return nil, util.ErrNotFound
}

func (f *fakeHanfuRepo) List(_ context.Context, dynasty, size, form string, page, pageSize int) ([]model.Hanfu, int64, error) {
	var out []model.Hanfu
	for _, h := range f.hanfus {
		if dynasty != "" && h.Dynasty != dynasty {
			continue
		}
		out = append(out, *h)
	}
	return out, int64(len(out)), nil
}

func (f *fakeHanfuRepo) UpdateStock(_ context.Context, id uint, stock int) error {
	existing, err := f.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	_ = existing
	f.hanfus[id].Stock = stock
	return nil
}

func (f *fakeHanfuRepo) DecrementStock(_ context.Context, id uint) error {
	existing, err := f.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	if existing.Stock <= 0 {
		return util.ErrConflict
	}
	f.hanfus[id].Stock--
	return nil
}

func (f *fakeHanfuRepo) IncrementStock(_ context.Context, id uint) error {
	_, err := f.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	f.hanfus[id].Stock++
	return nil
}

func (f *fakeHanfuRepo) Count(context.Context) (int64, error) { return int64(len(f.hanfus)), nil }

func TestHanfuServiceCreate(t *testing.T) {
	svc := NewHanfuService(newFakeHanfuRepo(), slog.Default())
	tests := []struct {
		name    string
		dynasty string
		wantErr bool
	}{
		{name: "valid tang", dynasty: constants.HanfuDynastyTang, wantErr: false},
		{name: "valid ming", dynasty: constants.HanfuDynastyMing, wantErr: false},
		{name: "invalid dynasty", dynasty: "qing", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &dto.CreateHanfuRequest{Name: "测试汉服", Dynasty: tt.dynasty, Size: "M", PricePerDay: 100, Stock: 2}
			_, err := svc.Create(context.Background(), req)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestHanfuServiceUpdateStockNotFound(t *testing.T) {
	svc := NewHanfuService(newFakeHanfuRepo(), slog.Default())
	_, err := svc.UpdateStock(context.Background(), 999, 5)
	if err == nil {
		t.Fatalf("expected not found error")
	}
}
