package result

import (
	"testing"

	"golang.org/x/text/message"
)

func TestEvaluateExperiment(t *testing.T) {
	tests := []struct {
		n        int
		pValue   float64
		expected string
	}{
		{n: 30, pValue: 0.15, expected: "Sample size too small to draw reliable conclusions. More data is needed."},
		{n: 30, pValue: 0.07, expected: "Results are marginally significant, but the small sample size limits reliability. Collect more data."},
		{n: 30, pValue: 0.03, expected: "Statistical significance reached, but the small sample size limits confidence. Validation with more data is recommended."},
		{n: 70, pValue: 0.15, expected: "Results are not significant. A larger sample size may help detect subtle effects."},
		{n: 70, pValue: 0.07, expected: "Results are marginally significant. Consider increasing sample size for validation."},
		{n: 70, pValue: 0.03, expected: "Results are statistically significant, but a larger sample size would strengthen confidence."},
		{n: 150, pValue: 0.15, expected: "Results are not significant. More participants may improve statistical power."},
		{n: 150, pValue: 0.07, expected: "Results are marginally significant. Consider more data to confirm findings."},
		{n: 150, pValue: 0.03, expected: "Results are statistically significant, supported by an adequate sample size."},
		{n: 350, pValue: 0.15, expected: "Results are not significant, even with a large sample size. Effect may be too small or non-existent."},
		{n: 350, pValue: 0.07, expected: "Results are marginally significant, suggesting a possible effect. Further analysis recommended."},
		{n: 350, pValue: 0.03, expected: "Results are statistically significant and supported by a large sample size, providing robust evidence."},
	}

	printer := message.NewPrinter(message.MatchLanguage("en"))

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := EvaluateExperiment(tt.n, tt.pValue, printer)
			if result != tt.expected {
				t.Errorf("EvaluateExperiment(%d, %f) = %v; want %v", tt.n, tt.pValue, result, tt.expected)
			}
		})
	}
}
