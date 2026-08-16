package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/dto"
	"github.com/lp/hanfu-rental/internal/util"
)

func TestHanfuLifecycleValidation(t *testing.T) {
	svc := NewHanfuService(newFakeHanfuRepo(), slog.Default())
	_, err := svc.Create(context.Background(), &dto.CreateHanfuRequest{Name: "清朝汉服", Dynasty: "qing", Size: "M", PricePerDay: 100, Stock: 2})
	if err == nil {
		t.Fatal("invalid dynasty should be rejected")
	}
	if util.DynastyText(constants.HanfuDynastyTang) != "唐制" {
		t.Fatalf("DynastyText(tang) = %s, want 唐制", util.DynastyText(constants.HanfuDynastyTang))
	}
}

func TestHanfuUpdateStock(t *testing.T) {
	svc := NewHanfuService(newFakeHanfuRepo(), slog.Default())
	created, _ := svc.Create(context.Background(), &dto.CreateHanfuRequest{Name: "唐制襦裙", Dynasty: constants.HanfuDynastyTang, Size: "M", PricePerDay: 100, Stock: 2})
	updated, err := svc.UpdateStock(context.Background(), created.ID, 5)
	if err != nil {
		t.Fatalf("update stock failed: %v", err)
	}
	if updated.Stock != 5 {
		t.Fatalf("stock = %d, want 5", updated.Stock)
	}
}
