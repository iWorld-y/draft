package biz

import (
	"backend/pkg/translator"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewGreeterUsecase,
	NewDictionaryUseCase,
	NewLearningUseCase,
	NewAuthUseCase,
	ProvideTranslator,
)

// ProvideTranslator 提供翻译器
func ProvideTranslator() translator.Translator {
	// 这里应该读取配置文件，暂时返回 nil
	return nil
}
