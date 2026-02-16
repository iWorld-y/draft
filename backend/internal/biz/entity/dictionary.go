// internal/biz/entity/dictionary.go
package entity

import (
	"time"
)

// Dictionary 词典实体
type Dictionary struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	TotalWords   int       `json:"total_words" db:"total_words"`
	LearnedWords int       `json:"learned_words" db:"learned_words"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Progress 计算学习进度
func (d *Dictionary) Progress() float64 {
	if d.TotalWords == 0 {
		return 0
	}
	return float64(d.LearnedWords) / float64(d.TotalWords) * 100
}

// Word 单词实体
type Word struct {
	ID             int64                  `json:"id" db:"id"`
	DictID         int64                  `json:"dict_id" db:"dict_id"`
	Word           string                 `json:"word" db:"word"`
	Phonetic       string                 `json:"phonetic" db:"phonetic"`
	Meaning        map[string]interface{} `json:"meaning" db:"meaning"`
	Example        string                 `json:"example" db:"example"`
	AudioURL       string                 `json:"audio_url" db:"audio_url"`
	Status         string                 `json:"status" db:"status"`           // new/learning/review/mastered
	EFFactor       float64                `json:"ef_factor" db:"ef_factor"`     // 遗忘因子
	Interval       int                    `json:"interval" db:"interval"`       // 间隔天数
	Repetitions    int                    `json:"repetitions" db:"repetitions"` // 已复习次数
	NextReviewDate *time.Time             `json:"next_review_date" db:"next_review_date"`
	LastReviewDate *time.Time             `json:"last_review_date" db:"last_review_date"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" db:"updated_at"`
}

// LearnRecord 学习记录实体
type LearnRecord struct {
	ID             int64     `json:"id" db:"id"`
	WordID         int64     `json:"word_id" db:"word_id"`
	Quality        int       `json:"quality" db:"quality"`
	TimeSpent      int       `json:"time_spent" db:"time_spent"`
	EFFactorBefore float64   `json:"ef_factor_before" db:"ef_factor_before"`
	EFFactorAfter  float64   `json:"ef_factor_after" db:"ef_factor_after"`
	IntervalBefore int       `json:"interval_before" db:"interval_before"`
	IntervalAfter  int       `json:"interval_after" db:"interval_after"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// UploadTask 上传任务实体
type UploadTask struct {
	ID             string     `json:"id" db:"id"`
	DictID         *int64     `json:"dict_id" db:"dict_id"`
	Status         string     `json:"status" db:"status"` // pending/processing/completed/failed
	TotalWords     int        `json:"total_words" db:"total_words"`
	ProcessedWords int        `json:"processed_words" db:"processed_words"`
	FailedWords    []string   `json:"failed_words" db:"failed_words"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	CompletedAt    *time.Time `json:"completed_at" db:"completed_at"`
}

// Progress 计算进度百分比
func (t *UploadTask) Progress() float64 {
	if t.TotalWords == 0 {
		return 0
	}
	return float64(t.ProcessedWords) / float64(t.TotalWords) * 100
}
