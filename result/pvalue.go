package result

import (
	"golang.org/x/text/message"
)

func EvaluateExperiment(n int, pValue float64, printer *message.Printer) string {
	switch {
	case n < 50:
		if pValue >= 0.10 {
			return printer.Sprintf("Sample size too small to draw reliable conclusions. More data is needed.")
		} else if pValue >= 0.05 && pValue < 0.10 {
			return printer.Sprintf("Results are marginally significant, but the small sample size limits reliability. Collect more data.")
		} else {
			return printer.Sprintf("Statistical significance reached, but the small sample size limits confidence. Validation with more data is recommended.")
		}
	case n >= 50 && n < 100:
		if pValue >= 0.10 {
			return printer.Sprintf("Results are not significant. A larger sample size may help detect subtle effects.")
		} else if pValue >= 0.05 && pValue < 0.10 {
			return printer.Sprintf("Results are marginally significant. Consider increasing sample size for validation.")
		} else {
			return printer.Sprintf("Results are statistically significant, but a larger sample size would strengthen confidence.")
		}
	case n >= 100 && n < 300:
		if pValue >= 0.10 {
			return printer.Sprintf("Results are not significant. More participants may improve statistical power.")
		} else if pValue >= 0.05 && pValue < 0.10 {
			return printer.Sprintf("Results are marginally significant. Consider more data to confirm findings.")
		} else {
			return printer.Sprintf("Results are statistically significant, supported by an adequate sample size.")
		}
	default: // N >= 300
		if pValue >= 0.10 {
			return printer.Sprintf("Results are not significant, even with a large sample size. Effect may be too small or non-existent.")
		} else if pValue >= 0.05 && pValue < 0.10 {
			return printer.Sprintf("Results are marginally significant, suggesting a possible effect. Further analysis recommended.")
		} else {
			return printer.Sprintf("Results are statistically significant and supported by a large sample size, providing robust evidence.")
		}
	}
}
