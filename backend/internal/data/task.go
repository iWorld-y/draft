// internal/data/task.go
package data

import (
	"context"
	"encoding/json"
	"time"

	"backend/internal/biz/entity"
	"backend/internal/biz/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type uploadTaskRepo struct {
	data *Data
	log  *log.Helper
}

// NewUploadTaskRepo 创建上传任务仓库实例
func NewUploadTaskRepo(data *Data, logger log.Logger) repo.UploadTaskRepo {
	return &uploadTaskRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Create 创建任务
func (r *uploadTaskRepo) Create(ctx context.Context, task *entity.UploadTask) error {
	query := `
		INSERT INTO upload_tasks (id, dict_id, status, total_words, processed_words, failed_words, failed_details, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	failedJSON, _ := json.Marshal(task.FailedWords)
	failedDetailsJSON, _ := json.Marshal(task.FailedDetails)
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	_, err := r.data.db.ExecContext(ctx, query,
		task.ID, task.DictID, task.Status, task.TotalWords,
		task.ProcessedWords, failedJSON, failedDetailsJSON, task.CreatedAt, task.UpdatedAt,
	)
	if err != nil {
		r.log.Errorf("failed to create upload task: %v", err)
		return err
	}
	return nil
}

// GetByID 根据 ID 获取任务
func (r *uploadTaskRepo) GetByID(ctx context.Context, id string) (*entity.UploadTask, error) {
	query := `
		SELECT id, dict_id, status, total_words, processed_words, failed_words, failed_details, created_at, updated_at, completed_at
		FROM upload_tasks
		WHERE id = $1
	`
	task := &entity.UploadTask{}
	var failedJSON []byte
	var failedDetailsJSON []byte
	var completedAt *time.Time
	err := r.data.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.DictID, &task.Status, &task.TotalWords,
		&task.ProcessedWords, &failedJSON, &failedDetailsJSON, &task.CreatedAt, &task.UpdatedAt, &completedAt,
	)
	if err != nil {
		r.log.Errorf("failed to get upload task: %v", err)
		return nil, err
	}
	task.CompletedAt = completedAt
	json.Unmarshal(failedJSON, &task.FailedWords)
	json.Unmarshal(failedDetailsJSON, &task.FailedDetails)
	return task, nil
}

// Update 更新任务
func (r *uploadTaskRepo) Update(ctx context.Context, task *entity.UploadTask) error {
	query := `
		UPDATE upload_tasks
		SET status = $1, processed_words = $2, failed_words = $3, failed_details = $4, updated_at = $5, completed_at = $6
		WHERE id = $7
	`
	failedJSON, _ := json.Marshal(task.FailedWords)
	failedDetailsJSON, _ := json.Marshal(task.FailedDetails)
	task.UpdatedAt = time.Now()

	_, err := r.data.db.ExecContext(ctx, query,
		task.Status, task.ProcessedWords, failedJSON, failedDetailsJSON,
		task.UpdatedAt, task.CompletedAt, task.ID,
	)
	if err != nil {
		r.log.Errorf("failed to update upload task: %v", err)
		return err
	}
	return nil
}

// IncrementProcessed 增加已处理数量
func (r *uploadTaskRepo) IncrementProcessed(ctx context.Context, id string, count int) error {
	query := `
		UPDATE upload_tasks
		SET processed_words = processed_words + $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.data.db.ExecContext(ctx, query, count, time.Now(), id)
	if err != nil {
		r.log.Errorf("failed to increment processed: %v", err)
		return err
	}
	return nil
}

// AddFailedWord 添加失败单词
func (r *uploadTaskRepo) AddFailedWord(ctx context.Context, id string, word string) error {
	query := `
		UPDATE upload_tasks
		SET failed_words = failed_words || $1::jsonb, updated_at = $2
		WHERE id = $3
	`
	wordJSON, _ := json.Marshal([]string{word})
	_, err := r.data.db.ExecContext(ctx, query, wordJSON, time.Now(), id)
	if err != nil {
		r.log.Errorf("failed to add failed word: %v", err)
		return err
	}
	return nil
}

// AddFailedWordWithReason 添加失败单词和失败原因
func (r *uploadTaskRepo) AddFailedWordWithReason(ctx context.Context, id, word, stage, reason string) error {
	query := `
		UPDATE upload_tasks
		SET failed_words = failed_words || $1::jsonb,
		    failed_details = failed_details || $2::jsonb,
		    updated_at = $3
		WHERE id = $4
	`

	wordJSON, _ := json.Marshal([]string{word})
	detailJSON, _ := json.Marshal([]entity.FailedDetail{
		{
			Word:   word,
			Stage:  stage,
			Reason: reason,
			At:     time.Now(),
		},
	})
	_, err := r.data.db.ExecContext(ctx, query, wordJSON, detailJSON, time.Now(), id)
	if err != nil {
		r.log.Errorf("failed to add failed word with reason: %v", err)
		return err
	}
	return nil
}
