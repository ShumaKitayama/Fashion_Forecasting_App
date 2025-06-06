// backend/internal/prediction/arima.go
package prediction

// PredictARIMA は時系列データに対するARIMAモデルによる予測を行います
// 注: 完全なARIMA実装は複雑なため、簡易版として線形予測を返します
func PredictARIMA(data []float64, horizon int) ([]float64, error) {
	// 簡易実装として線形回帰を使用
	return PredictLinearRegression(data, horizon)
}

// PredictLinearRegression は線形回帰による予測を行います
func PredictLinearRegression(data []float64, horizon int) ([]float64, error) {
	if len(data) < 2 {
		return nil, nil
	}
	
	// 単純な線形回帰のパラメータを計算
	n := float64(len(data))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumXX := 0.0
	
	for i, y := range data {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}
	
	// 傾きと切片を計算
	slope := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)
	intercept := (sumY - slope*sumX) / n
	
	// 予測値を計算
	predictions := make([]float64, horizon)
	for i := 0; i < horizon; i++ {
		x := float64(len(data) + i)
		predictions[i] = slope*x + intercept
		
		// 負の値を防止
		if predictions[i] < 0 {
			predictions[i] = 0
		}
	}
	
	return predictions, nil
}