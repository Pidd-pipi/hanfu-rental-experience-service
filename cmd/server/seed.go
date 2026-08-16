package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// seed populates default users, hanfu, makeup packages, photographers, activities and articles.
func seed(ctx context.Context, db *gorm.DB, logger *slog.Logger) error {
	var total int64
	if err := db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return fmt.Errorf("seed count users: %w", err)
	}
	if total > 0 {
		logger.Info("seed skipped: users already exist", slog.Int64("count", total))
		return nil
	}
	hash := func(p string) (string, error) {
		b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	type seedUser struct {
		phone, password, nickname, role string
		deposit                         float64
	}
	defs := []seedUser{
		{phone: "13500000001", password: "123456", nickname: "汉服小主", role: constants.UserRoleCustomer, deposit: 500},
		{phone: "13500000002", password: "123456", nickname: "同袍阿琳", role: constants.UserRoleCustomer, deposit: 0},
		{phone: "13600000001", password: "admin123", nickname: "掌柜", role: constants.UserRoleAdmin, deposit: 0},
	}
	users := make([]model.User, 0, len(defs))
	for _, d := range defs {
		h, err := hash(d.password)
		if err != nil {
			return fmt.Errorf("seed hash user %s: %w", d.phone, err)
		}
		users = append(users, model.User{
			Phone: d.phone, PasswordHash: h, Nickname: d.nickname,
			Role: d.role, Deposit: d.deposit,
		})
	}
	if err := db.WithContext(ctx).Create(&users).Error; err != nil {
		return fmt.Errorf("seed users: %w", err)
	}
	hanfus := []model.Hanfu{
		{Name: "齐胸襦裙·绯云", Dynasty: constants.HanfuDynastyTang, Form: "齐胸襦裙", Color: "绯红", Size: "S/M", PricePerDay: 128, Images: "", Stock: 3, Status: constants.HanfuStatusAvailable},
		{Name: "宋制褙子·青芜", Dynasty: constants.HanfuDynastySong, Form: "褙子", Color: "青绿", Size: "M/L", PricePerDay: 98, Images: "", Stock: 4, Status: constants.HanfuStatusAvailable},
		{Name: "明制马面裙·黛蓝", Dynasty: constants.HanfuDynastyMing, Form: "马面裙", Color: "黛蓝", Size: "M", PricePerDay: 108, Images: "", Stock: 2, Status: constants.HanfuStatusAvailable},
		{Name: "汉制曲裾·素白", Dynasty: constants.HanfuDynastyHan, Form: "曲裾", Color: "素白", Size: "S/M", PricePerDay: 118, Images: "", Stock: 2, Status: constants.HanfuStatusAvailable},
	}
	if err := db.WithContext(ctx).Create(&hanfus).Error; err != nil {
		return fmt.Errorf("seed hanfus: %w", err)
	}
	makeups := []model.MakeupPackage{
		{Name: "基础妆造", Description: "底妆+眉形+口红", Price: 88},
		{Name: "精致妆造", Description: "含发型与配饰", Price: 168},
		{Name: "复原妆造", Description: "朝代复原发型妆容", Price: 268},
	}
	if err := db.WithContext(ctx).Create(&makeups).Error; err != nil {
		return fmt.Errorf("seed makeup packages: %w", err)
	}
	photographers := []model.Photographer{
		{Name: "阿黎", Style: "古风人像", PricePerDay: 500},
		{Name: "小鹤", Style: "园林外景", PricePerDay: 680},
	}
	if err := db.WithContext(ctx).Create(&photographers).Error; err != nil {
		return fmt.Errorf("seed photographers: %w", err)
	}
	activities := []model.Activity{
		{Title: "上巳节雅集", Category: "雅集", StartTime: time.Now().Add(72 * time.Hour), Location: "南湖公园", Flow: "签到-焚香-品茶-踏青", DressCode: "汉服或常服", Fee: 0, MaxParticipants: 30, Status: "open"},
		{Title: "汉服走秀展", Category: "汉服走秀", StartTime: time.Now().Add(120 * time.Hour), Location: "文化馆", Flow: "彩排-走秀-评奖", DressCode: "汉服", Fee: 50, MaxParticipants: 20, Status: "open"},
	}
	if err := db.WithContext(ctx).Create(&activities).Error; err != nil {
		return fmt.Errorf("seed activities: %w", err)
	}
	articles := []model.Article{
		{Title: "汉服形制入门：襦裙与褙子", Category: "形制科普", Content: "汉服按朝代与形制可分为汉制、唐制、宋制、明制……本文介绍常见的襦裙、褙子、马面裙等形制的特点与穿着要点。", CoverURL: "", PublishedAt: time.Now()},
		{Title: "穿汉服出门的穿搭指南", Category: "穿搭指南", Content: "从发型、配饰到鞋子，教你如何把汉服穿出日常感，让传统服饰融入现代生活。", CoverURL: "", PublishedAt: time.Now()},
	}
	if err := db.WithContext(ctx).Create(&articles).Error; err != nil {
		return fmt.Errorf("seed articles: %w", err)
	}
	logger.Info(fmt.Sprintf(constants.LogSeedingCompleted, len(users), len(hanfus)))
	return nil
}
