// internal/service/dictionary.go
package service

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"backend/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// DictionaryService 词典服务
type DictionaryService struct {
	uc  *biz.DictionaryUseCase
	log *log.Helper
}

// NewDictionaryService 创建词典服务
func NewDictionaryService(uc *biz.DictionaryUseCase, logger log.Logger) *DictionaryService {
	return &DictionaryService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// CreateDictionaryRequest 创建词典请求
type CreateDictionaryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateDictionaryResponse 创建词典响应
type CreateDictionaryResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TotalWords  int     `json:"total_words"`
	Progress    float64 `json:"progress"`
}

// CreateDictionary 创建词典
func (s *DictionaryService) CreateDictionary(ctx context.Context, req *CreateDictionaryRequest) (*CreateDictionaryResponse, error) {
	// 暂时硬编码用户 ID 为 1
	dict, err := s.uc.CreateDictionary(ctx, req.Name, req.Description, 1)
	if err != nil {
		return nil, err
	}

	return &CreateDictionaryResponse{
		ID:          dict.ID,
		Name:        dict.Name,
		Description: dict.Description,
		TotalWords:  dict.TotalWords,
		Progress:    dict.Progress(),
	}, nil
}

// ListDictionariesResponse 词典列表响应
type ListDictionariesResponse struct {
	Items []*DictionaryItem `json:"items"`
}

// DictionaryItem 词典项
type DictionaryItem struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	TotalWords   int     `json:"total_words"`
	LearnedWords int     `json:"learned_words"`
	Progress     float64 `json:"progress"`
	CreatedAt    string  `json:"created_at"`
}

// ListDictionaries 获取词典列表
func (s *DictionaryService) ListDictionaries(ctx context.Context) (*ListDictionariesResponse, error) {
	dicts, err := s.uc.ListDictionaries(ctx, 1)
	if err != nil {
		return nil, err
	}

	items := make([]*DictionaryItem, 0, len(dicts))
	for _, dict := range dicts {
		items = append(items, &DictionaryItem{
			ID:           dict.ID,
			Name:         dict.Name,
			Description:  dict.Description,
			TotalWords:   dict.TotalWords,
			LearnedWords: dict.LearnedWords,
			Progress:     dict.Progress(),
			CreatedAt:    dict.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &ListDictionariesResponse{Items: items}, nil
}

// UploadDictionaryRequest 上传词典请求
type UploadDictionaryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UploadDictionaryResponse 上传词典响应
type UploadDictionaryResponse struct {
	TaskID         string `json:"task_id"`
	Status         string `json:"status"`
	TotalWords     int    `json:"total_words"`
	ProcessedWords int    `json:"processed_words"`
}

// UploadDictionary 上传词典文件
func (s *DictionaryService) UploadDictionary(ctx context.Context) (*UploadDictionaryResponse, error) {
	// 从 HTTP 请求中获取文件
	httpCtx, ok := http.RequestFromServerContext(ctx)
	if !ok {
		return nil, fmt.Errorf("not http context")
	}

	// 解析 multipart form
	if err := httpCtx.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		return nil, fmt.Errorf("failed to parse form: %w", err)
	}

	file, _, err := httpCtx.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	defer file.Close()

	name := httpCtx.FormValue("name")
	if name == "" {
		name = "未命名词典"
	}
	description := httpCtx.FormValue("description")

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	result, err := s.uc.UploadDictionary(ctx, bytes.NewReader(content), name, description, 1)
	if err != nil {
		return nil, err
	}

	return &UploadDictionaryResponse{
		TaskID:         result.TaskID,
		Status:         result.Status,
		TotalWords:     result.TotalWords,
		ProcessedWords: result.ProcessedWords,
	}, nil
}

// GetUploadStatusRequest 获取上传状态请求
type GetUploadStatusRequest struct {
	TaskID string `json:"task_id"`
}

// GetUploadStatusResponse 获取上传状态响应
type GetUploadStatusResponse struct {
	TaskID      string   `json:"task_id"`
	Status      string   `json:"status"`
	Progress    float64  `json:"progress"`
	Total       int      `json:"total"`
	Processed   int      `json:"processed"`
	FailedWords []string `json:"failed_words"`
}

// GetUploadStatus 获取上传任务状态
func (s *DictionaryService) GetUploadStatus(ctx context.Context, req *GetUploadStatusRequest) (*GetUploadStatusResponse, error) {
	task, err := s.uc.GetUploadStatus(ctx, req.TaskID)
	if err != nil {
		return nil, err
	}

	return &GetUploadStatusResponse{
		TaskID:      task.ID,
		Status:      task.Status,
		Progress:    task.Progress(),
		Total:       task.TotalWords,
		Processed:   task.ProcessedWords,
		FailedWords: task.FailedWords,
	}, nil
}
