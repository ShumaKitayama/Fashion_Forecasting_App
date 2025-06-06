package prediction

import "testing"

func TestPredictLinearRegressionLength(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	horizon := 3

	preds, err := PredictLinearRegression(data, horizon)
	if err != nil {
		t.Fatalf("PredictLinearRegression returned error: %v", err)
	}
	if len(preds) != horizon {
		t.Fatalf("expected %d predictions, got %d", horizon, len(preds))
	}
}

func TestPredictLinearRegressionNonNegative(t *testing.T) {
	data := []float64{5, 4, 3, 2, 1}
	horizon := 3

	preds, err := PredictLinearRegression(data, horizon)
	if err != nil {
		t.Fatalf("PredictLinearRegression returned error: %v", err)
	}
	for i, p := range preds {
		if p < 0 {
			t.Fatalf("prediction %d is negative: %f", i, p)
		}
	}
}
