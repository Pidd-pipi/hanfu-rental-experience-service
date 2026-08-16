package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/dto"
	"github.com/lp/hanfu-rental/internal/model"
	"github.com/lp/hanfu-rental/internal/repository"
	"github.com/lp/hanfu-rental/internal/util"
)

// ActivityService manages cultural activities and signups.
type ActivityService struct {
	activities *repository.ActivityRepository
	logger     *slog.Logger
}

// NewActivityService wires the activity service dependencies.
func NewActivityService(activities *repository.ActivityRepository, logger *slog.Logger) *ActivityService {
	return &ActivityService{activities: activities, logger: logger}
}

// Create publishes a new activity.
func (s *ActivityService) Create(ctx context.Context, req *dto.CreateActivityRequest) (*model.Activity, error) {
	a := &model.Activity{
		Title: req.Title, Category: req.Category, StartTime: req.StartTime,
		Location: req.Location, Flow: req.Flow, DressCode: req.DressCode,
		Fee: req.Fee, MaxParticipants: req.MaxParticipants, Status: "open",
	}
	if err := s.activities.Create(ctx, a); err != nil {
		return nil, util.WrapAppError(fmt.Errorf("activity[title=%s] create: %w", req.Title, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogActivityCreateSuccess, a.ID, a.Title))
	return a, nil
}

// List returns activities with pagination.
func (s *ActivityService) List(ctx context.Context, q *dto.PageQuery) (*dto.PageResult, error) {
	q.Normalize()
	items, total, err := s.activities.List(ctx, q.Page, q.PageSize)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("activity list: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

// Get returns one activity.
func (s *ActivityService) Get(ctx context.Context, id uint) (*model.Activity, error) {
	a, err := s.activities.FindByID(ctx, id)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("activity[id=%d] get: %w", id, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	return a, nil
}

// Signup registers a customer for an activity. The activity row is locked
// inside a transaction so capacity and duplicate checks are race-free, and a
// unique (activity_id,user_id) index guards against concurrent duplicates.
func (s *ActivityService) Signup(ctx context.Context, user *model.User, activityID uint) (*model.ActivityRegistration, error) {
	reg := &model.ActivityRegistration{}
	if err := s.activities.Transaction(ctx, func(txCtx context.Context) error {
		activity, err := s.activities.FindByIDForUpdate(txCtx, activityID)
		if err != nil {
			if errors.Is(err, util.ErrNotFound) {
				return util.WrapAppError(fmt.Errorf("activity[id=%d] signup find: %w", activityID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
			}
			return err
		}
		count, err := s.activities.CountRegistrations(txCtx, activityID)
		if err != nil {
			return err
		}
		if count >= int64(activity.MaxParticipants) {
			return util.NewAppError(409, constants.CodeConflict, constants.MsgActivityClosed, nil)
		}
		if _, err := s.activities.FindRegistration(txCtx, activityID, user.ID); err == nil {
			return util.NewAppError(409, constants.CodeConflict, constants.MsgAlreadySignedUp, nil)
		} else if !errors.Is(err, util.ErrNotFound) {
			return err
		}
		reg = &model.ActivityRegistration{ActivityID: activityID, UserID: user.ID, Fee: activity.Fee}
		if err := s.activities.CreateRegistration(txCtx, reg); err != nil {
			return err
		}
		return nil
	}); err != nil {
		var appErr *util.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		if errors.Is(err, util.ErrNotFound) {
			return nil, util.WrapAppError(fmt.Errorf("activity[id=%d] signup find: %w", activityID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
		}
		s.logger.Error(fmt.Sprintf(constants.LogActivitySignupFailed, activityID, user.ID, err))
		return nil, util.WrapAppError(fmt.Errorf("activity[id=%d] signup: %w", activityID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogActivitySignupSuccess, activityID, user.ID))
	return reg, nil
}
