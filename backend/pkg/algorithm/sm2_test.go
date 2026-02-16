// pkg/algorithm/sm2_test.go
package algorithm

import (
	"testing"
	"time"
)

func TestCalculateNextReview(t *testing.T) {
	tests := []struct {
		name            string
		quality         int
		currentEF       float64
		currentInterval int
		repetitions     int
		wantInterval    int
		wantEFMin       float64
		wantEFMax       float64
	}{
		{
			name:            "首次学习答对",
			quality:         4,
			currentEF:       2.5,
			currentInterval: 0,
			repetitions:     0,
			wantInterval:    1,
			wantEFMin:       2.5,
			wantEFMax:       2.7,
		},
		{
			name:            "第二次复习答对",
			quality:         5,
			currentEF:       2.5,
			currentInterval: 1,
			repetitions:     1,
			wantInterval:    6,
			wantEFMin:       2.5, // quality=5 时, EF 增加 0.1: 2.5 + 0.1 = 2.6, 但受到 max 2.5 限制
			wantEFMax:       2.5,
		},
		{
			name:            "答错重置",
			quality:         1,
			currentEF:       2.5,
			currentInterval: 6,
			repetitions:     2,
			wantInterval:    1,
			wantEFMin:       1.96, // 2.5 + (0.1 - 4*0.16) = 2.5 - 0.54 = 1.96
			wantEFMax:       1.96,
		},
		{
			name:            "第三次复习答对",
			quality:         4,
			currentEF:       2.5,
			currentInterval: 6,
			repetitions:     2,
			wantInterval:    15,  // 6 * 2.5 = 15
			wantEFMin:       2.5, // 2.5 + (0.1 - 1*0.1) = 2.5
			wantEFMax:       2.5,
		},
		{
			name:            "完全记住",
			quality:         5,
			currentEF:       2.4,
			currentInterval: 1,
			repetitions:     1,
			wantInterval:    6,
			wantEFMin:       2.5,
			wantEFMax:       2.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateNextReview(
				tt.quality,
				tt.currentEF,
				tt.currentInterval,
				tt.repetitions,
			)

			if result.Interval != tt.wantInterval {
				t.Errorf("CalculateNextReview() Interval = %v, want %v", result.Interval, tt.wantInterval)
			}
			if result.EFactor < tt.wantEFMin || result.EFactor > tt.wantEFMax {
				t.Errorf("CalculateNextReview() EFactor = %v, want between %v and %v", result.EFactor, tt.wantEFMin, tt.wantEFMax)
			}
			if result.NextReviewDate.IsZero() {
				t.Error("CalculateNextReview() NextReviewDate should not be zero")
			}
		})
	}
}

func TestCalculateNextReview_EFactorBounds(t *testing.T) {
	// 测试 EF 下限
	result := CalculateNextReview(0, 1.3, 10, 5)
	if result.EFactor < 1.3 {
		t.Errorf("EFactor should not be less than 1.3, got %v", result.EFactor)
	}

	// 测试 EF 上限
	result = CalculateNextReview(5, 2.5, 10, 5)
	if result.EFactor > 2.5 {
		t.Errorf("EFactor should not be greater than 2.5, got %v", result.EFactor)
	}
}

func TestGetQualityDescription(t *testing.T) {
	tests := []struct {
		quality  int
		wantDesc string
	}{
		{0, "完全不认识"},
		{1, "有印象但想不起来"},
		{2, "想起来了但很费力"},
		{3, "有些犹豫但想起来了"},
		{4, "轻松想起来"},
		{5, "脱口而出"},
		{6, "未知"},
		{-1, "未知"},
	}

	for _, tt := range tests {
		t.Run(tt.wantDesc, func(t *testing.T) {
			got := GetQualityDescription(tt.quality)
			if got != tt.wantDesc {
				t.Errorf("GetQualityDescription(%d) = %v, want %v", tt.quality, got, tt.wantDesc)
			}
		})
	}
}

func TestGetWordStatus(t *testing.T) {
	tests := []struct {
		interval    int
		repetitions int
		wantStatus  string
	}{
		{0, 0, "new"},
		{1, 0, "learning"},
		{6, 1, "review"},
		{30, 3, "mastered"},
		{60, 5, "mastered"},
	}

	for _, tt := range tests {
		t.Run(tt.wantStatus, func(t *testing.T) {
			got := GetWordStatus(tt.interval, tt.repetitions)
			if got != tt.wantStatus {
				t.Errorf("GetWordStatus(%d, %d) = %v, want %v", tt.interval, tt.repetitions, got, tt.wantStatus)
			}
		})
	}
}

func TestCalculateNextReview_NextReviewDate(t *testing.T) {
	result := CalculateNextReview(4, 2.5, 0, 0)

	expectedDate := time.Now().AddDate(0, 0, result.Interval)
	// 允许 1 秒误差
	diff := result.NextReviewDate.Sub(expectedDate)
	if diff < -time.Second || diff > time.Second {
		t.Errorf("NextReviewDate = %v, want approximately %v", result.NextReviewDate, expectedDate)
	}
}
