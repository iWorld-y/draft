// pkg/algorithm/sm2.go
package algorithm

import (
	"math"
	"time"
)

// SM2Result SM-2 算法计算结果
type SM2Result struct {
	EFactor        float64   // 新的遗忘因子
	Interval       int       // 新的间隔天数
	Repetitions    int       // 新的复习次数
	NextReviewDate time.Time // 下次复习日期
}

// CalculateNextReview SM-2 算法核心函数
// quality: 答题质量 0-5
// currentEF: 当前 E-Factor
// currentInterval: 当前间隔天数
// repetitions: 已复习次数
func CalculateNextReview(
	quality int,
	currentEF float64,
	currentInterval int,
	repetitions int,
) SM2Result {
	// 1. 计算新的 E-Factor
	// 公式: EF' = EF + (0.1 - (5-q) * (0.08 + (5-q) * 0.02))
	newEF := currentEF + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))

	// E-Factor 范围限制在 1.3 - 2.5
	if newEF < 1.3 {
		newEF = 1.3
	}
	if newEF > 2.5 {
		newEF = 2.5
	}

	var newInterval int
	var newRepetitions int

	// 2. 根据答题质量决定间隔
	if quality < 3 {
		// 答错了，重新开始
		newInterval = 1
		newRepetitions = 0
	} else {
		newRepetitions = repetitions + 1

		// 根据复习次数计算间隔
		switch newRepetitions {
		case 1:
			newInterval = 1
		case 2:
			newInterval = 6
		default:
			// 使用公式: I(n) = I(n-1) * EF
			newInterval = int(math.Round(float64(currentInterval) * newEF))
		}
	}

	// 3. 计算下次复习日期
	nextReviewDate := time.Now().AddDate(0, 0, newInterval)

	return SM2Result{
		EFactor:        newEF,
		Interval:       newInterval,
		Repetitions:    newRepetitions,
		NextReviewDate: nextReviewDate,
	}
}

// GetQualityDescription 获取质量等级描述
func GetQualityDescription(quality int) string {
	descriptions := map[int]string{
		0: "完全不认识",
		1: "有印象但想不起来",
		2: "想起来了但很费力",
		3: "有些犹豫但想起来了",
		4: "轻松想起来",
		5: "脱口而出",
	}
	if desc, ok := descriptions[quality]; ok {
		return desc
	}
	return "未知"
}

// GetWordStatus 根据复习参数获取单词状态
func GetWordStatus(interval int, repetitions int) string {
	if interval >= 30 {
		return "mastered" // 已掌握
	} else if repetitions > 0 {
		return "review" // 复习中
	} else if interval > 0 {
		return "learning" // 学习中
	}
	return "new" // 新词
}
