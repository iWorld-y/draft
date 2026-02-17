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
	return translator.NewFreeDictionaryTranslator("")
}
