package service

import (
	"bytes"
	"context"

	v1 "backend/api/helloworld/v1"
	authctx "backend/internal/auth"
	"backend/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// DictionaryService 词典服务
type DictionaryService struct {
	v1.UnimplementedDictionaryServer

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

// CreateDictionary 创建词典
func (s *DictionaryService) CreateDictionary(ctx context.Context, req *v1.CreateDictionaryRequest) (*v1.CreateDictionaryReply, error) {
	userID, ok := authctx.UserIDFromContext(ctx)
	if !ok || userID <= 0 {
		return nil, biz.ErrUnauthorized
	}
	dict, err := s.uc.CreateDictionary(ctx, req.Name, req.Description, userID)
	if err != nil {
		return nil, err
	}

	return &v1.CreateDictionaryReply{
		Id:          dict.ID,
		Name:        dict.Name,
		Description: dict.Description,
		TotalWords:  int32(dict.TotalWords),
		Progress:    dict.Progress(),
	}, nil
}

// ListDictionaries 获取词典列表
func (s *DictionaryService) ListDictionaries(ctx context.Context, _ *v1.ListDictionariesRequest) (*v1.ListDictionariesReply, error) {
	userID, ok := authctx.UserIDFromContext(ctx)
	if !ok || userID <= 0 {
		return nil, biz.ErrUnauthorized
	}
	dicts, err := s.uc.ListDictionaries(ctx, userID)
	if err != nil {
		return nil, err
	}

	items := make([]*v1.DictionaryItem, 0, len(dicts))
	for _, dict := range dicts {
		items = append(items, &v1.DictionaryItem{
			Id:           dict.ID,
			Name:         dict.Name,
			Description:  dict.Description,
			TotalWords:   int32(dict.TotalWords),
			LearnedWords: int32(dict.LearnedWords),
			Progress:     dict.Progress(),
			CreatedAt:    dict.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &v1.ListDictionariesReply{Items: items}, nil
}

// UploadDictionary 上传词典文件
func (s *DictionaryService) UploadDictionary(ctx context.Context, req *v1.UploadDictionaryRequest) (*v1.UploadDictionaryReply, error) {
	userID, ok := authctx.UserIDFromContext(ctx)
	if !ok || userID <= 0 {
		return nil, biz.ErrUnauthorized
	}

	name := req.Name
	if name == "" {
		name = "未命名词典"
	}

	result, err := s.uc.UploadDictionary(ctx, bytes.NewReader(req.FileContent), name, req.Description, userID)
	if err != nil {
		s.log.Warnf("upload dictionary failed, user_id=%d name=%q: %v", userID, name, err)
		return nil, err
	}

	return &v1.UploadDictionaryReply{
		TaskId:         result.TaskID,
		Status:         result.Status,
		TotalWords:     int32(result.TotalWords),
		ProcessedWords: int32(result.ProcessedWords),
	}, nil
}

// GetUploadStatus 获取上传任务状态
func (s *DictionaryService) GetUploadStatus(ctx context.Context, req *v1.GetUploadStatusRequest) (*v1.GetUploadStatusReply, error) {
	userID, ok := authctx.UserIDFromContext(ctx)
	if !ok || userID <= 0 {
		return nil, biz.ErrUnauthorized
	}
	task, err := s.uc.GetUploadStatus(ctx, req.TaskId, userID)
	if err != nil {
		return nil, err
	}

	return &v1.GetUploadStatusReply{
		TaskId:      task.ID,
		Status:      task.Status,
		Progress:    task.Progress(),
		Total:       int32(task.TotalWords),
		Processed:   int32(task.ProcessedWords),
		FailedWords: task.FailedWords,
	}, nil
}
