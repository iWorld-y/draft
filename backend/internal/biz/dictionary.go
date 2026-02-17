// internal/biz/dictionary.go
package biz

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"backend/internal/biz/entity"
	"backend/internal/biz/repo"
	"backend/pkg/translator"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrEmptyWordFile = kerrors.BadRequest("EMPTY_WORD_FILE", "文件中没有可导入的单词")
)

// DictionaryUseCase 词典业务逻辑
type DictionaryUseCase struct {
	dictRepo   repo.DictionaryRepo
	wordRepo   repo.WordRepo
	taskRepo   repo.UploadTaskRepo
	translator translator.Translator
	log        *log.Helper
}

// NewDictionaryUseCase 创建词典业务逻辑实例
func NewDictionaryUseCase(
	dictRepo repo.DictionaryRepo,
	wordRepo repo.WordRepo,
	taskRepo repo.UploadTaskRepo,
	translator translator.Translator,
	logger log.Logger,
) *DictionaryUseCase {
	return &DictionaryUseCase{
		dictRepo:   dictRepo,
		wordRepo:   wordRepo,
		taskRepo:   taskRepo,
		translator: translator,
		log:        log.NewHelper(logger),
	}
}

// CreateDictionary 创建词典
func (uc *DictionaryUseCase) CreateDictionary(ctx context.Context, name, description string, userID int64) (*entity.Dictionary, error) {
	dict := &entity.Dictionary{
		UserID:      userID,
		Name:        name,
		Description: description,
	}
	if err := uc.dictRepo.Create(ctx, dict); err != nil {
		return nil, fmt.Errorf("failed to create dictionary: %w", err)
	}
	return dict, nil
}

// GetDictionary 获取词典详情
func (uc *DictionaryUseCase) GetDictionary(ctx context.Context, id int64) (*entity.Dictionary, error) {
	return uc.dictRepo.GetByID(ctx, id)
}

// ListDictionaries 获取词典列表
func (uc *DictionaryUseCase) ListDictionaries(ctx context.Context, userID int64) ([]*entity.Dictionary, error) {
	return uc.dictRepo.ListByUserID(ctx, userID)
}

// UploadTaskResult 上传任务结果
type UploadTaskResult struct {
	TaskID         string `json:"task_id"`
	Status         string `json:"status"`
	TotalWords     int    `json:"total_words"`
	ProcessedWords int    `json:"processed_words"`
}

// UploadDictionary 上传词典文件
func (uc *DictionaryUseCase) UploadDictionary(ctx context.Context, reader io.Reader, name, description string, userID int64) (*UploadTaskResult, error) {
	// 1. 解析文件，提取单词列表
	words, err := uc.parseWordFile(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse word file: %w", err)
	}

	// 2. 创建词典记录
	dict, err := uc.CreateDictionary(ctx, name, description, userID)
	if err != nil {
		return nil, err
	}

	// 3. 创建上传任务
	taskID := fmt.Sprintf("task_%d_%d", dict.ID, time.Now().Unix())
	task := &entity.UploadTask{
		ID:            taskID,
		DictID:        &dict.ID,
		Status:        "processing",
		TotalWords:    len(words),
		FailedWords:   []string{},
		FailedDetails: []entity.FailedDetail{},
	}
	if err := uc.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to create upload task: %w", err)
	}

	// 4. 启动异步任务处理
	go uc.processUploadTask(taskID, dict.ID, userID, words)

	return &UploadTaskResult{
		TaskID:         taskID,
		Status:         "processing",
		TotalWords:     len(words),
		ProcessedWords: 0,
	}, nil
}

// parseWordFile 解析单词文件
func (uc *DictionaryUseCase) parseWordFile(reader io.Reader) ([]string, error) {
	var words []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			words = append(words, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(words) == 0 {
		return nil, ErrEmptyWordFile
	}
	return words, nil
}

// processUploadTask 异步处理上传任务
func (uc *DictionaryUseCase) processUploadTask(taskID string, dictID, userID int64, words []string) {
	ctx := context.Background()
	total := len(words)

	// 并发控制：每次最多 5 个并发
	semaphore := make(chan struct{}, 5)
	done := make(chan bool, total)

	for _, wordStr := range words {
		semaphore <- struct{}{} // 获取信号量

		go func(w string) {
			defer func() { <-semaphore }() // 释放信号量

			// 检查是否已存在
			existing, _ := uc.wordRepo.GetByDictIDAndWord(ctx, dictID, w)
			if existing != nil {
				// 已存在，跳过
				uc.taskRepo.IncrementProcessed(ctx, taskID, 1)
				done <- true
				return
			}

			// 跨词典复用：若该用户库内已有该词，直接复用释义并跳过 API 请求
			cachedWord, _ := uc.wordRepo.GetByUserAndWord(ctx, userID, w)
			if cachedWord != nil {
				word := &entity.Word{
					DictID:   dictID,
					Word:     w,
					Phonetic: cachedWord.Phonetic,
					Meaning:  cachedWord.Meaning,
					Example:  cachedWord.Example,
					AudioURL: cachedWord.AudioURL,
					Status:   "new",
				}
				if err := uc.wordRepo.Create(ctx, word); err != nil {
					uc.recordUploadFailure(ctx, taskID, w, "reuse", err)
				}
				uc.taskRepo.IncrementProcessed(ctx, taskID, 1)
				done <- true
				return
			}

			// 调用翻译 API
			detail, err := uc.translator.Translate(w)
			if err != nil {
				// 翻译失败，记录失败单词
				uc.recordUploadFailure(ctx, taskID, w, "translate", err)
				uc.taskRepo.IncrementProcessed(ctx, taskID, 1)
				done <- true
				return
			}

			// 保存到数据库
			word := &entity.Word{
				DictID:   dictID,
				Word:     detail.Word,
				Phonetic: detail.Phonetic,
				Meaning:  detail.Meaning,
				Example:  detail.Example,
				Status:   "new",
			}
			if err := uc.wordRepo.Create(ctx, word); err != nil {
				uc.recordUploadFailure(ctx, taskID, w, "save", err)
			}

			// 更新进度
			uc.taskRepo.IncrementProcessed(ctx, taskID, 1)

			// 速率限制：防止 API 限流
			time.Sleep(100 * time.Millisecond)
			done <- true
		}(wordStr)
	}

	// 等待所有任务完成
	for i := 0; i < total; i++ {
		<-done
	}

	// 更新任务状态为完成
	task, _ := uc.taskRepo.GetByID(ctx, taskID)
	if task != nil {
		// 若全部处理都失败，则标记任务失败，避免前端误判“成功”
		if total > 0 && len(task.FailedWords) >= total {
			task.Status = "failed"
		} else {
			task.Status = "completed"
		}
		now := time.Now()
		task.CompletedAt = &now
		uc.taskRepo.Update(ctx, task)

		// 更新词典统计
		count, _ := uc.wordRepo.CountByDictID(ctx, dictID)
		uc.dictRepo.UpdateStats(ctx, dictID, count, 0)
	}
}

func truncateReason(reason string) string {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return "unknown error"
	}

	const maxReasonLength = 512
	if len(reason) <= maxReasonLength {
		return reason
	}
	return reason[:maxReasonLength]
}

func (uc *DictionaryUseCase) recordUploadFailure(ctx context.Context, taskID, word, stage string, err error) {
	reason := "unknown error"
	if err != nil {
		reason = truncateReason(err.Error())
	}

	uc.log.WithContext(ctx).Warnf("Upload task word failed task_id=%s word=%q stage=%s reason=%q", taskID, word, stage, reason)
	if saveErr := uc.taskRepo.AddFailedWordWithReason(ctx, taskID, word, stage, reason); saveErr != nil {
		uc.log.WithContext(ctx).Errorf("Failed to persist upload failure detail task_id=%s word=%q stage=%s err=%v", taskID, word, stage, saveErr)
	}
}

// GetUploadStatus 获取上传任务状态
func (uc *DictionaryUseCase) GetUploadStatus(ctx context.Context, taskID string, userID int64) (*entity.UploadTask, error) {
	task, err := uc.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil || task.DictID == nil {
		return nil, ErrUnauthorized
	}

	dict, err := uc.dictRepo.GetByID(ctx, *task.DictID)
	if err != nil {
		return nil, err
	}
	if dict == nil || dict.UserID != userID {
		return nil, ErrUnauthorized
	}
	return task, nil
}
