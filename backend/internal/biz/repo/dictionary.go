// internal/biz/repo/dictionary.go
package repo

import (
	"context"

	"backend/internal/biz/entity"
)

// DictionaryRepo 词典仓库接口
type DictionaryRepo interface {
	// Create 创建词典
	Create(ctx context.Context, dict *entity.Dictionary) error
	// GetByID 根据 ID 获取词典
	GetByID(ctx context.Context, id int64) (*entity.Dictionary, error)
	// ListByUserID 获取用户的词典列表
	ListByUserID(ctx context.Context, userID int64) ([]*entity.Dictionary, error)
	// Update 更新词典
	Update(ctx context.Context, dict *entity.Dictionary) error
	// Delete 删除词典（软删除）
	Delete(ctx context.Context, id int64) error
	// UpdateStats 更新词典统计信息
	UpdateStats(ctx context.Context, id int64, totalWords, learnedWords int) error
	// IsOwnedByUser 判断词典是否属于该用户
	IsOwnedByUser(ctx context.Context, dictID, userID int64) (bool, error)
}

// WordRepo 单词仓库接口
type WordRepo interface {
	// Create 创建单词
	Create(ctx context.Context, word *entity.Word) error
	// CreateBatch 批量创建单词
	CreateBatch(ctx context.Context, words []*entity.Word) error
	// GetByID 根据 ID 获取单词
	GetByID(ctx context.Context, id int64) (*entity.Word, error)
	// GetByIDForUser 根据用户归属获取单词
	GetByIDForUser(ctx context.Context, id, userID int64) (*entity.Word, error)
	// GetByDictIDAndWord 根据词典 ID 和单词获取
	GetByDictIDAndWord(ctx context.Context, dictID int64, word string) (*entity.Word, error)
	// GetByUserAndWord 根据用户和单词获取（跨词典复用）
	GetByUserAndWord(ctx context.Context, userID int64, word string) (*entity.Word, error)
	// ListByDictID 获取词典的单词列表
	ListByDictID(ctx context.Context, dictID int64, offset, limit int) ([]*entity.Word, error)
	// CountByDictID 统计词典单词数
	CountByDictID(ctx context.Context, dictID int64) (int, error)
	// Update 更新单词
	Update(ctx context.Context, word *entity.Word) error
	// GetTodayTasks 获取今日学习任务
	GetTodayTasks(ctx context.Context, dictID int64, limit int) ([]*entity.Word, error)
	// CountReviewToday 统计今日待复习数
	CountReviewToday(ctx context.Context, dictID int64) (int, error)
	// CountNewWords 统计新词数
	CountNewWords(ctx context.Context, dictID int64) (int, error)
}

// LearnRecordRepo 学习记录仓库接口
type LearnRecordRepo interface {
	// Create 创建学习记录
	Create(ctx context.Context, record *entity.LearnRecord) error
	// ListByWordID 获取单词的学习记录
	ListByWordID(ctx context.Context, wordID int64, limit int) ([]*entity.LearnRecord, error)
}

// UploadTaskRepo 上传任务仓库接口
type UploadTaskRepo interface {
	// Create 创建任务
	Create(ctx context.Context, task *entity.UploadTask) error
	// GetByID 根据 ID 获取任务
	GetByID(ctx context.Context, id string) (*entity.UploadTask, error)
	// Update 更新任务
	Update(ctx context.Context, task *entity.UploadTask) error
	// IncrementProcessed 增加已处理数量
	IncrementProcessed(ctx context.Context, id string, count int) error
	// AddFailedWord 添加失败单词
	AddFailedWord(ctx context.Context, id string, word string) error
}
