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

// ArticleService manages hanfu culture science articles.
type ArticleService struct {
	articles *repository.ArticleRepository
	logger   *slog.Logger
}

// NewArticleService wires the article service dependencies.
func NewArticleService(articles *repository.ArticleRepository, logger *slog.Logger) *ArticleService {
	return &ArticleService{articles: articles, logger: logger}
}

// Create publishes an article.
func (s *ArticleService) Create(ctx context.Context, req *dto.CreateArticleRequest) (*model.Article, error) {
	a := &model.Article{
		Title: req.Title, Category: req.Category, Content: req.Content,
		CoverURL: req.CoverURL, PublishedAt: time.Now(),
	}
	if err := s.articles.Create(ctx, a); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogArticleCreateFailed, req.Title, err))
		return nil, util.WrapAppError(fmt.Errorf("article[title=%s] create: %w", req.Title, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogArticleCreateSuccess, a.ID, a.Title))
	return a, nil
}

// List returns articles with pagination.
func (s *ArticleService) List(ctx context.Context, q *dto.PageQuery) (*dto.PageResult, error) {
	q.Normalize()
	items, total, err := s.articles.List(ctx, q.Page, q.PageSize)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("article list: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

// Get returns one article.
func (s *ArticleService) Get(ctx context.Context, id uint) (*model.Article, error) {
	a, err := s.articles.FindByID(ctx, id)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("article[id=%d] get: %w", id, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	return a, nil
}
