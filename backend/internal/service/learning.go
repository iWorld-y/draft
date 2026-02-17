package service

import (
	"context"
	"encoding/json"

	v1 "backend/api/helloworld/v1"
	authctx "backend/internal/auth"
	"backend/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// LearningService 学习服务
type LearningService struct {
	v1.UnimplementedLearningServer

	uc  *biz.LearningUseCase
	log *log.Helper
}

// NewLearningService 创建学习服务
func NewLearningService(uc *biz.LearningUseCase, logger log.Logger) *LearningService {
	return &LearningService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// GetTodayTasks 获取今日学习任务
func (s *LearningService) GetTodayTasks(ctx context.Context, req *v1.GetTodayTasksRequest) (*v1.GetTodayTasksReply, error) {
	userID, ok := authctx.UserIDFromContext(ctx)
	if !ok || userID <= 0 {
		return nil, biz.ErrUnauthorized
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 20
	}

	result, err := s.uc.GetTodayTasks(ctx, userID, req.DictId, limit)
	if err != nil {
		return nil, err
	}

	words := make([]*v1.WordItem, 0, len(result.Words))
	for _, w := range result.Words {
		meaningJSON, _ := json.Marshal(w.Meaning)
		nextReview := ""
		if w.NextReviewDate != nil {
			nextReview = w.NextReviewDate.Format("2006-01-02")
		}
		words = append(words, &v1.WordItem{
			Id:             w.ID,
			Word:           w.Word,
			Phonetic:       w.Phonetic,
			Meaning:        meaningJSON,
			Example:        w.Example,
			AudioUrl:       w.AudioURL,
			Status:         w.Status,
			NextReviewDate: nextReview,
		})
	}

	return &v1.GetTodayTasksReply{
		ReviewCount: int32(result.ReviewCount),
		NewCount:    int32(result.NewCount),
		Words:       words,
	}, nil
}

// SubmitLearning 提交学习结果
func (s *LearningService) SubmitLearning(ctx context.Context, req *v1.SubmitLearningRequest) (*v1.SubmitLearningReply, error) {
	userID, ok := authctx.UserIDFromContext(ctx)
	if !ok || userID <= 0 {
		return nil, biz.ErrUnauthorized
	}

	quality := int(req.Quality)
	if quality < 0 || quality > 5 {
		return nil, biz.ErrInvalidInput
	}

	result, err := s.uc.SubmitLearning(ctx, userID, req.WordId, quality, int(req.TimeSpent))
	if err != nil {
		return nil, err
	}

	return &v1.SubmitLearningReply{
		WordId:         result.WordID,
		NewStatus:      result.NewStatus,
		NewInterval:    int32(result.NewInterval),
		NextReviewDate: result.NextReviewDate.Format("2006-01-02"),
		EfFactor:       result.EFFactor,
	}, nil
}
