// internal/service/learning.go
package service

import (
	"context"
	"encoding/json"
	"strconv"

	"backend/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// LearningService 学习服务
type LearningService struct {
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

// TodayTasksRequest 今日学习任务请求
type TodayTasksRequest struct {
	DictID string `json:"dict_id"`
	Limit  string `json:"limit"`
}

// Meaning 释义结构
type Meaning struct {
	Definitions []Definition `json:"definitions"`
}

// Definition 释义定义
type Definition struct {
	Pos  string `json:"pos"`
	Text string `json:"text"`
}

// WordItem 单词项
type WordItem struct {
	ID             int64           `json:"id"`
	Word           string          `json:"word"`
	Phonetic       string          `json:"phonetic"`
	Meaning        json.RawMessage `json:"meaning"`
	Example        string          `json:"example"`
	AudioURL       string          `json:"audio_url"`
	Status         string          `json:"status"`
	NextReviewDate string          `json:"next_review_date"`
}

// TodayTasksResponse 今日学习任务响应
type TodayTasksResponse struct {
	ReviewCount int         `json:"review_count"`
	NewCount    int         `json:"new_count"`
	Words       []*WordItem `json:"words"`
}

// GetTodayTasks 获取今日学习任务
func (s *LearningService) GetTodayTasks(ctx context.Context, req *TodayTasksRequest) (*TodayTasksResponse, error) {
	dictID, err := strconv.ParseInt(req.DictID, 10, 64)
	if err != nil {
		return nil, err
	}

	limit := 20
	if req.Limit != "" {
		if l, err := strconv.Atoi(req.Limit); err == nil && l > 0 {
			limit = l
		}
	}

	result, err := s.uc.GetTodayTasks(ctx, dictID, limit)
	if err != nil {
		return nil, err
	}

	words := make([]*WordItem, 0, len(result.Words))
	for _, w := range result.Words {
		meaningJSON, _ := json.Marshal(w.Meaning)
		nextReview := ""
		if w.NextReviewDate != nil {
			nextReview = w.NextReviewDate.Format("2006-01-02")
		}
		words = append(words, &WordItem{
			ID:             w.ID,
			Word:           w.Word,
			Phonetic:       w.Phonetic,
			Meaning:        meaningJSON,
			Example:        w.Example,
			AudioURL:       w.AudioURL,
			Status:         w.Status,
			NextReviewDate: nextReview,
		})
	}

	return &TodayTasksResponse{
		ReviewCount: result.ReviewCount,
		NewCount:    result.NewCount,
		Words:       words,
	}, nil
}

// SubmitLearningRequest 提交学习请求
type SubmitLearningRequest struct {
	WordID    string `json:"word_id"`
	Quality   string `json:"quality"`
	TimeSpent string `json:"time_spent"`
}

// SubmitLearningResponse 提交学习响应
type SubmitLearningResponse struct {
	WordID         int64   `json:"word_id"`
	NewStatus      string  `json:"new_status"`
	NewInterval    int     `json:"new_interval"`
	NextReviewDate string  `json:"next_review_date"`
	EFFactor       float64 `json:"ef_factor"`
}

// SubmitLearning 提交学习结果
func (s *LearningService) SubmitLearning(ctx context.Context, req *SubmitLearningRequest) (*SubmitLearningResponse, error) {
	wordID, err := strconv.ParseInt(req.WordID, 10, 64)
	if err != nil {
		return nil, err
	}

	quality, err := strconv.Atoi(req.Quality)
	if err != nil || quality < 0 || quality > 5 {
		return nil, err
	}

	timeSpent := 0
	if req.TimeSpent != "" {
		timeSpent, _ = strconv.Atoi(req.TimeSpent)
	}

	result, err := s.uc.SubmitLearning(ctx, wordID, quality, timeSpent)
	if err != nil {
		return nil, err
	}

	return &SubmitLearningResponse{
		WordID:         result.WordID,
		NewStatus:      result.NewStatus,
		NewInterval:    result.NewInterval,
		NextReviewDate: result.NextReviewDate.Format("2006-01-02"),
		EFFactor:       result.EFFactor,
	}, nil
}
