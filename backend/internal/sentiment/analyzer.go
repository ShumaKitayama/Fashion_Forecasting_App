// backend/internal/sentiment/analyzer.go
package sentiment

import (
	"strings"
)

// 感情分析器
type Analyzer struct {
}

// 新しい分析器の作成
func NewAnalyzer() (*Analyzer, error) {
	return &Analyzer{}, nil
}

// テキストの感情分析を実行
func (a *Analyzer) AnalyzeText(text string) (float64, error) {
	// テキスト前処理
	text = strings.ToLower(text)
	
	// 簡易実装：キーワードベースの感情分析
	positiveWords := []string{"love", "great", "amazing", "beautiful", "perfect", "excellent", "stylish", "trendy"}
	negativeWords := []string{"hate", "bad", "poor", "terrible", "ugly", "worst", "outdated", "boring"}
	
	var positiveScore, negativeScore float64
	
	// ポジティブワードのカウント
	for _, word := range positiveWords {
		if strings.Contains(text, word) {
			positiveScore += 0.2
		}
	}
	
	// ネガティブワードのカウント
	for _, word := range negativeWords {
		if strings.Contains(text, word) {
			negativeScore += 0.2
		}
	}
	
	// スコアの正規化（-1.0から1.0の範囲に）
	score := positiveScore - negativeScore
	if score > 1.0 {
		score = 1.0
	} else if score < -1.0 {
		score = -1.0
	}
	
	return score, nil
}

// 複数テキストの感情分析を実行
func (a *Analyzer) AnalyzeTexts(texts []string) (SentimentResult, error) {
	var positive, neutral, negative float64
	
	for _, text := range texts {
		score, err := a.AnalyzeText(text)
		if err != nil {
			continue
		}
		
		if score > 0.3 {
			positive++
		} else if score < -0.3 {
			negative++
		} else {
			neutral++
		}
	}
	
	total := positive + neutral + negative
	if total == 0 {
		return SentimentResult{0.33, 0.34, 0.33}, nil
	}
	
	return SentimentResult{
		Positive: positive / total,
		Neutral:  neutral / total,
		Negative: negative / total,
	}, nil
}

// 感情分析結果構造体
type SentimentResult struct {
	Positive float64 `json:"positive"`
	Neutral  float64 `json:"neutral"`
	Negative float64 `json:"negative"`
}