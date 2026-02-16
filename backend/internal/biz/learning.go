// internal/biz/learning.go
package biz

import (
	"context"
	"fmt"
	"time"

	"backend/internal/biz/entity"
	"backend/internal/biz/repo"
	"backend/pkg/algorithm"
)

// LearningUseCase 学习业务逻辑
type LearningUseCase struct {
	wordRepo   repo.WordRepo
	recordRepo repo.LearnRecordRepo
	dictRepo   repo.DictionaryRepo
}

// NewLearningUseCase 创建学习业务逻辑实例
func NewLearningUseCase(
	wordRepo repo.WordRepo,
	recordRepo repo.LearnRecordRepo,
	dictRepo repo.DictionaryRepo,
) *LearningUseCase {
	return &LearningUseCase{
		wordRepo:   wordRepo,
		recordRepo: recordRepo,
		dictRepo:   dictRepo,
	}
}

// TodayTasksResult 今日学习任务结果
type TodayTasksResult struct {
	ReviewCount int            `json:"review_count"`
	NewCount    int            `json:"new_count"`
	Words       []*entity.Word `json:"words"`
}

// GetTodayTasks 获取今日学习任务
func (uc *LearningUseCase) GetTodayTasks(ctx context.Context, dictID int64, limit int) (*TodayTasksResult, error) {
	// 1. 获取今日待复习数
	reviewCount, err := uc.wordRepo.CountReviewToday(ctx, dictID)
	if err != nil {
		return nil, fmt.Errorf("failed to count review tasks: %w", err)
	}

	// 2. 获取新词数
	newCount, err := uc.wordRepo.CountNewWords(ctx, dictID)
	if err != nil {
		return nil, fmt.Errorf("failed to count new words: %w", err)
	}

	// 3. 获取任务队列
	words, err := uc.wordRepo.GetTodayTasks(ctx, dictID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get today tasks: %w", err)
	}

	return &TodayTasksResult{
		ReviewCount: reviewCount,
		NewCount:    newCount,
		Words:       words,
	}, nil
}

// SubmitResult 提交学习结果
type SubmitResult struct {
	WordID         int64     `json:"word_id"`
	NewStatus      string    `json:"new_status"`
	NewInterval    int       `json:"new_interval"`
	NextReviewDate time.Time `json:"next_review_date"`
	EFFactor       float64   `json:"ef_factor"`
}

// SubmitLearning 提交学习结果
func (uc *LearningUseCase) SubmitLearning(ctx context.Context, wordID int64, quality, timeSpent int) (*SubmitResult, error) {
	// 1. 查询单词当前状态
	word, err := uc.wordRepo.GetByID(ctx, wordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get word: %w", err)
	}

	// 2. 记录学习前的状态
	oldEF := word.EFFactor
	oldInterval := word.Interval

	// 3. 调用 SM-2 算法计算新参数
	result := algorithm.CalculateNextReview(
		quality,
		word.EFFactor,
		word.Interval,
		word.Repetitions,
	)

	// 4. 更新单词记忆参数
	word.EFFactor = result.EFactor
	word.Interval = result.Interval
	word.Repetitions = result.Repetitions
	word.NextReviewDate = &result.NextReviewDate
	now := time.Now()
	word.LastReviewDate = &now

	// 更新状态
	word.Status = algorithm.GetWordStatus(result.Interval, result.Repetitions)

	// 5. 保存更新
	if err := uc.wordRepo.Update(ctx, word); err != nil {
		return nil, fmt.Errorf("failed to update word: %w", err)
	}

	// 6. 记录学习日志
	record := &entity.LearnRecord{
		WordID:         wordID,
		Quality:        quality,
		TimeSpent:      timeSpent,
		EFFactorBefore: oldEF,
		EFFactorAfter:  result.EFactor,
		IntervalBefore: oldInterval,
		IntervalAfter:  result.Interval,
	}
	if err := uc.recordRepo.Create(ctx, record); err != nil {
		return nil, fmt.Errorf("failed to create learn record: %w", err)
	}

	// 7. 更新词典统计（如果单词首次变为已学习状态）
	if word.Status != "new" && oldInterval == 0 {
		uc.updateDictionaryStats(ctx, word.DictID)
	}

	return &SubmitResult{
		WordID:         wordID,
		NewStatus:      word.Status,
		NewInterval:    result.Interval,
		NextReviewDate: result.NextReviewDate,
		EFFactor:       result.EFactor,
	}, nil
}

// updateDictionaryStats 更新词典统计
func (uc *LearningUseCase) updateDictionaryStats(ctx context.Context, dictID int64) {
	// 这里可以添加更新词典学习进度的逻辑
	// 例如：增加 learned_words 计数
}
