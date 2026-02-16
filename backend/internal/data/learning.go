// internal/data/learning.go
package data

import (
	"context"
	"time"

	"backend/internal/biz/entity"
	"backend/internal/biz/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type learnRecordRepo struct {
	data *Data
	log  *log.Helper
}

// NewLearnRecordRepo 创建学习记录仓库实例
func NewLearnRecordRepo(data *Data, logger log.Logger) repo.LearnRecordRepo {
	return &learnRecordRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Create 创建学习记录
func (r *learnRecordRepo) Create(ctx context.Context, record *entity.LearnRecord) error {
	query := `
		INSERT INTO learn_records (word_id, quality, time_spent, ef_factor_before, ef_factor_after, interval_before, interval_after, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	record.CreatedAt = time.Now()

	err := r.data.db.QueryRowContext(ctx, query,
		record.WordID, record.Quality, record.TimeSpent,
		record.EFFactorBefore, record.EFFactorAfter,
		record.IntervalBefore, record.IntervalAfter,
		record.CreatedAt,
	).Scan(&record.ID)

	if err != nil {
		r.log.Errorf("failed to create learn record: %v", err)
		return err
	}
	return nil
}

// ListByWordID 获取单词的学习记录
func (r *learnRecordRepo) ListByWordID(ctx context.Context, wordID int64, limit int) ([]*entity.LearnRecord, error) {
	query := `
		SELECT id, word_id, quality, time_spent, ef_factor_before, ef_factor_after, interval_before, interval_after, created_at
		FROM learn_records
		WHERE word_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.data.db.QueryContext(ctx, query, wordID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*entity.LearnRecord
	for rows.Next() {
		record := &entity.LearnRecord{}
		err := rows.Scan(
			&record.ID, &record.WordID, &record.Quality, &record.TimeSpent,
			&record.EFFactorBefore, &record.EFFactorAfter,
			&record.IntervalBefore, &record.IntervalAfter,
			&record.CreatedAt,
		)
		if err != nil {
			continue
		}
		records = append(records, record)
	}
	return records, nil
}
