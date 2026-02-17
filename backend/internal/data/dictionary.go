// internal/data/dictionary.go
package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"backend/internal/biz/entity"
	"backend/internal/biz/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type dictionaryRepo struct {
	data *Data
	log  *log.Helper
}

// NewDictionaryRepo 创建词典仓库实例
func NewDictionaryRepo(data *Data, logger log.Logger) repo.DictionaryRepo {
	return &dictionaryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Create 创建词典
func (r *dictionaryRepo) Create(ctx context.Context, dict *entity.Dictionary) error {
	query := `
		INSERT INTO dictionaries (user_id, name, description, total_words, learned_words, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	now := time.Now()
	dict.CreatedAt = now
	dict.UpdatedAt = now

	err := r.data.db.QueryRowContext(ctx, query,
		dict.UserID, dict.Name, dict.Description,
		dict.TotalWords, dict.LearnedWords,
		dict.CreatedAt, dict.UpdatedAt,
	).Scan(&dict.ID)

	if err != nil {
		r.log.Errorf("failed to create dictionary: %v", err)
		return err
	}
	return nil
}

// GetByID 根据 ID 获取词典
func (r *dictionaryRepo) GetByID(ctx context.Context, id int64) (*entity.Dictionary, error) {
	query := `
		SELECT id, user_id, name, description, total_words, learned_words, created_at, updated_at
		FROM dictionaries
		WHERE id = $1 AND deleted_at IS NULL
	`
	dict := &entity.Dictionary{}
	err := r.data.db.QueryRowContext(ctx, query, id).Scan(
		&dict.ID, &dict.UserID, &dict.Name, &dict.Description,
		&dict.TotalWords, &dict.LearnedWords,
		&dict.CreatedAt, &dict.UpdatedAt,
	)
	if err != nil {
		r.log.Errorf("failed to get dictionary: %v", err)
		return nil, err
	}
	return dict, nil
}

// ListByUserID 获取用户的词典列表
func (r *dictionaryRepo) ListByUserID(ctx context.Context, userID int64) ([]*entity.Dictionary, error) {
	query := `
		SELECT id, user_id, name, description, total_words, learned_words, created_at, updated_at
		FROM dictionaries
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`
	rows, err := r.data.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.log.Errorf("failed to list dictionaries: %v", err)
		return nil, err
	}
	defer rows.Close()

	var dicts []*entity.Dictionary
	for rows.Next() {
		dict := &entity.Dictionary{}
		err := rows.Scan(
			&dict.ID, &dict.UserID, &dict.Name, &dict.Description,
			&dict.TotalWords, &dict.LearnedWords,
			&dict.CreatedAt, &dict.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("failed to scan dictionary: %v", err)
			continue
		}
		dicts = append(dicts, dict)
	}
	return dicts, nil
}

// Update 更新词典
func (r *dictionaryRepo) Update(ctx context.Context, dict *entity.Dictionary) error {
	query := `
		UPDATE dictionaries
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
	`
	dict.UpdatedAt = time.Now()
	_, err := r.data.db.ExecContext(ctx, query,
		dict.Name, dict.Description, dict.UpdatedAt, dict.ID,
	)
	if err != nil {
		r.log.Errorf("failed to update dictionary: %v", err)
		return err
	}
	return nil
}

// Delete 删除词典（软删除）
func (r *dictionaryRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE dictionaries SET deleted_at = $1 WHERE id = $2`
	_, err := r.data.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		r.log.Errorf("failed to delete dictionary: %v", err)
		return err
	}
	return nil
}

// UpdateStats 更新词典统计信息
func (r *dictionaryRepo) UpdateStats(ctx context.Context, id int64, totalWords, learnedWords int) error {
	query := `
		UPDATE dictionaries
		SET total_words = $1, learned_words = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.data.db.ExecContext(ctx, query, totalWords, learnedWords, time.Now(), id)
	if err != nil {
		r.log.Errorf("failed to update dictionary stats: %v", err)
		return err
	}
	return nil
}

// IsOwnedByUser 判断词典是否属于该用户
func (r *dictionaryRepo) IsOwnedByUser(ctx context.Context, dictID, userID int64) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM dictionaries
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`
	var count int
	if err := r.data.db.QueryRowContext(ctx, query, dictID, userID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

// wordRepo 单词仓库实现
type wordRepo struct {
	data *Data
	log  *log.Helper
}

// NewWordRepo 创建单词仓库实例
func NewWordRepo(data *Data, logger log.Logger) repo.WordRepo {
	return &wordRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Create 创建单词
func (r *wordRepo) Create(ctx context.Context, word *entity.Word) error {
	query := `
		INSERT INTO words (dict_id, word, phonetic, meaning, example, audio_url, status, ef_factor, interval, repetitions, next_review_date, last_review_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`
	meaningJSON, _ := json.Marshal(word.Meaning)
	now := time.Now()
	word.CreatedAt = now
	word.UpdatedAt = now

	err := r.data.db.QueryRowContext(ctx, query,
		word.DictID, word.Word, word.Phonetic, meaningJSON, word.Example, word.AudioURL,
		word.Status, word.EFFactor, word.Interval, word.Repetitions,
		word.NextReviewDate, word.LastReviewDate,
		word.CreatedAt, word.UpdatedAt,
	).Scan(&word.ID)

	if err != nil {
		r.log.Errorf("failed to create word: %v", err)
		return err
	}
	return nil
}

// CreateBatch 批量创建单词
func (r *wordRepo) CreateBatch(ctx context.Context, words []*entity.Word) error {
	// 简化实现：逐个创建
	for _, word := range words {
		if err := r.Create(ctx, word); err != nil {
			r.log.Errorf("failed to create word in batch: %v", err)
			return err
		}
	}
	return nil
}

// GetByID 根据 ID 获取单词
func (r *wordRepo) GetByID(ctx context.Context, id int64) (*entity.Word, error) {
	query := `
		SELECT id, dict_id, word, phonetic, meaning, example, audio_url, status, ef_factor, interval, repetitions, next_review_date, last_review_date, created_at, updated_at
		FROM words
		WHERE id = $1
	`
	word := &entity.Word{}
	var meaningJSON []byte
	err := r.data.db.QueryRowContext(ctx, query, id).Scan(
		&word.ID, &word.DictID, &word.Word, &word.Phonetic, &meaningJSON, &word.Example,
		&word.AudioURL, &word.Status, &word.EFFactor, &word.Interval, &word.Repetitions,
		&word.NextReviewDate, &word.LastReviewDate, &word.CreatedAt, &word.UpdatedAt,
	)
	if err != nil {
		r.log.Errorf("failed to get word: %v", err)
		return nil, err
	}
	json.Unmarshal(meaningJSON, &word.Meaning)
	return word, nil
}

// GetByIDForUser 根据用户归属获取单词
func (r *wordRepo) GetByIDForUser(ctx context.Context, id, userID int64) (*entity.Word, error) {
	query := `
		SELECT w.id, w.dict_id, w.word, w.phonetic, w.meaning, w.example, w.audio_url, w.status, w.ef_factor, w.interval, w.repetitions, w.next_review_date, w.last_review_date, w.created_at, w.updated_at
		FROM words w
		INNER JOIN dictionaries d ON d.id = w.dict_id
		WHERE w.id = $1 AND d.user_id = $2 AND d.deleted_at IS NULL
	`
	word := &entity.Word{}
	var meaningJSON []byte
	err := r.data.db.QueryRowContext(ctx, query, id, userID).Scan(
		&word.ID, &word.DictID, &word.Word, &word.Phonetic, &meaningJSON, &word.Example,
		&word.AudioURL, &word.Status, &word.EFFactor, &word.Interval, &word.Repetitions,
		&word.NextReviewDate, &word.LastReviewDate, &word.CreatedAt, &word.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(meaningJSON, &word.Meaning)
	return word, nil
}

// GetByDictIDAndWord 根据词典 ID 和单词获取
func (r *wordRepo) GetByDictIDAndWord(ctx context.Context, dictID int64, wordStr string) (*entity.Word, error) {
	query := `
		SELECT id, dict_id, word, phonetic, meaning, example, audio_url, status, ef_factor, interval, repetitions, next_review_date, last_review_date, created_at, updated_at
		FROM words
		WHERE dict_id = $1 AND word = $2
	`
	word := &entity.Word{}
	var meaningJSON []byte
	err := r.data.db.QueryRowContext(ctx, query, dictID, wordStr).Scan(
		&word.ID, &word.DictID, &word.Word, &word.Phonetic, &meaningJSON, &word.Example,
		&word.AudioURL, &word.Status, &word.EFFactor, &word.Interval, &word.Repetitions,
		&word.NextReviewDate, &word.LastReviewDate, &word.CreatedAt, &word.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(meaningJSON, &word.Meaning)
	return word, nil
}

// GetByUserAndWord 根据用户和单词获取（跨词典复用）
func (r *wordRepo) GetByUserAndWord(ctx context.Context, userID int64, wordStr string) (*entity.Word, error) {
	query := `
		SELECT w.id, w.dict_id, w.word, w.phonetic, w.meaning, w.example, w.audio_url, w.status, w.ef_factor, w.interval, w.repetitions, w.next_review_date, w.last_review_date, w.created_at, w.updated_at
		FROM words w
		INNER JOIN dictionaries d ON d.id = w.dict_id
		WHERE d.user_id = $1 AND d.deleted_at IS NULL AND w.word = $2
		ORDER BY w.id ASC
		LIMIT 1
	`
	word := &entity.Word{}
	var meaningJSON []byte
	err := r.data.db.QueryRowContext(ctx, query, userID, wordStr).Scan(
		&word.ID, &word.DictID, &word.Word, &word.Phonetic, &meaningJSON, &word.Example,
		&word.AudioURL, &word.Status, &word.EFFactor, &word.Interval, &word.Repetitions,
		&word.NextReviewDate, &word.LastReviewDate, &word.CreatedAt, &word.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	json.Unmarshal(meaningJSON, &word.Meaning)
	return word, nil
}

// ListByDictID 获取词典的单词列表
func (r *wordRepo) ListByDictID(ctx context.Context, dictID int64, offset, limit int) ([]*entity.Word, error) {
	query := `
		SELECT id, dict_id, word, phonetic, meaning, example, audio_url, status, ef_factor, interval, repetitions, next_review_date, last_review_date, created_at, updated_at
		FROM words
		WHERE dict_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.data.db.QueryContext(ctx, query, dictID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*entity.Word
	for rows.Next() {
		word := &entity.Word{}
		var meaningJSON []byte
		err := rows.Scan(
			&word.ID, &word.DictID, &word.Word, &word.Phonetic, &meaningJSON, &word.Example,
			&word.AudioURL, &word.Status, &word.EFFactor, &word.Interval, &word.Repetitions,
			&word.NextReviewDate, &word.LastReviewDate, &word.CreatedAt, &word.UpdatedAt,
		)
		if err != nil {
			continue
		}
		json.Unmarshal(meaningJSON, &word.Meaning)
		words = append(words, word)
	}
	return words, nil
}

// CountByDictID 统计词典单词数
func (r *wordRepo) CountByDictID(ctx context.Context, dictID int64) (int, error) {
	query := `SELECT COUNT(*) FROM words WHERE dict_id = $1`
	var count int
	err := r.data.db.QueryRowContext(ctx, query, dictID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Update 更新单词
func (r *wordRepo) Update(ctx context.Context, word *entity.Word) error {
	query := `
		UPDATE words
		SET phonetic = $1, meaning = $2, example = $3, audio_url = $4, status = $5, ef_factor = $6, interval = $7, repetitions = $8, next_review_date = $9, last_review_date = $10, updated_at = $11
		WHERE id = $12
	`
	meaningJSON, _ := json.Marshal(word.Meaning)
	word.UpdatedAt = time.Now()

	_, err := r.data.db.ExecContext(ctx, query,
		word.Phonetic, meaningJSON, word.Example, word.AudioURL,
		word.Status, word.EFFactor, word.Interval, word.Repetitions,
		word.NextReviewDate, word.LastReviewDate, word.UpdatedAt, word.ID,
	)
	if err != nil {
		r.log.Errorf("failed to update word: %v", err)
		return err
	}
	return nil
}

// GetTodayTasks 获取今日学习任务
func (r *wordRepo) GetTodayTasks(ctx context.Context, dictID int64, limit int) ([]*entity.Word, error) {
	query := `
		SELECT id, dict_id, word, phonetic, meaning, example, audio_url, status, ef_factor, interval, repetitions, next_review_date, last_review_date, created_at, updated_at
		FROM words
		WHERE dict_id = $1
		AND (
			status = 'new'
			OR (next_review_date <= CURRENT_DATE AND status IN ('learning', 'review'))
		)
		ORDER BY 
			CASE WHEN next_review_date <= CURRENT_DATE THEN 0 ELSE 1 END,
			next_review_date ASC
		LIMIT $2
	`
	rows, err := r.data.db.QueryContext(ctx, query, dictID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*entity.Word
	for rows.Next() {
		word := &entity.Word{}
		var meaningJSON []byte
		err := rows.Scan(
			&word.ID, &word.DictID, &word.Word, &word.Phonetic, &meaningJSON, &word.Example,
			&word.AudioURL, &word.Status, &word.EFFactor, &word.Interval, &word.Repetitions,
			&word.NextReviewDate, &word.LastReviewDate, &word.CreatedAt, &word.UpdatedAt,
		)
		if err != nil {
			continue
		}
		json.Unmarshal(meaningJSON, &word.Meaning)
		words = append(words, word)
	}
	return words, nil
}

// CountReviewToday 统计今日待复习数
func (r *wordRepo) CountReviewToday(ctx context.Context, dictID int64) (int, error) {
	query := `
		SELECT COUNT(*) FROM words
		WHERE dict_id = $1
		AND next_review_date <= CURRENT_DATE
		AND status IN ('learning', 'review')
	`
	var count int
	err := r.data.db.QueryRowContext(ctx, query, dictID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountNewWords 统计新词数
func (r *wordRepo) CountNewWords(ctx context.Context, dictID int64) (int, error) {
	query := `SELECT COUNT(*) FROM words WHERE dict_id = $1 AND status = 'new'`
	var count int
	err := r.data.db.QueryRowContext(ctx, query, dictID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
