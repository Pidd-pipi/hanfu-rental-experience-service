package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/dto"
	"github.com/lp/hanfu-rental/internal/model"
	"github.com/lp/hanfu-rental/internal/repository"
	"github.com/lp/hanfu-rental/internal/util"
)

// MemberCardService manages member rental cards.
type MemberCardService struct {
	cards  *repository.MemberCardRepository
	logger *slog.Logger
}

// NewMemberCardService wires the member card service dependencies.
func NewMemberCardService(cards *repository.MemberCardRepository, logger *slog.Logger) *MemberCardService {
	return &MemberCardService{cards: cards, logger: logger}
}

// Create applies a member card for the customer.
func (s *MemberCardService) Create(ctx context.Context, userID uint, req *dto.CreateMemberCardRequest) (*model.MemberCard, error) {
	if _, err := s.cards.FindActiveByUser(ctx, userID); err == nil {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgCardExists, nil)
	}
	card := &model.MemberCard{
		UserID: userID, CardType: req.CardType, Status: constants.MemberCardStatusActive,
		ExpireAt: time.Now().AddDate(0, constants.MemberCardMonths(req.CardType), 0),
	}
	if err := s.cards.Create(ctx, card); err != nil {
		return nil, util.WrapAppError(fmt.Errorf("member_card[user=%d] create: %w", userID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogMemberCardCreateSuccess, card.ID, userID, req.CardType))
	return card, nil
}

// ListMy returns the cards of the current user.
func (s *MemberCardService) ListMy(ctx context.Context, userID uint) ([]model.MemberCard, error) {
	items, err := s.cards.ListByUser(ctx, userID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("member_card[user=%d] list: %w", userID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	return items, nil
}

// ProcessExpiry marks expired cards (called on startup or on demand).
func (s *MemberCardService) ProcessExpiry(ctx context.Context) (int64, error) {
	if err := s.cards.Expire(ctx, time.Now()); err != nil {
		return 0, util.WrapAppError(fmt.Errorf("member_card expire: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	return 0, nil
}
