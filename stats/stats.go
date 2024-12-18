package stats

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

// Data struct to hold the pre- and post-intervention scores
type Data struct {
	PreControl       float64
	PostControl      float64
	PreIntervention  float64
	PostIntervention float64
}

// Function to read data from CSV file
func ReadCSV(filePath string) ([]Data, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []Data
	for _, record := range records[1:] { // Skip header row
		preControl, _ := strconv.ParseFloat(record[0], 64)
		preIntervention, _ := strconv.ParseFloat(record[1], 64)
		postControl, _ := strconv.ParseFloat(record[2], 64)
		postIntervention, _ := strconv.ParseFloat(record[3], 64)

		data = append(data, Data{
			PreControl:       preControl,
			PostControl:      postControl,
			PreIntervention:  preIntervention,
			PostIntervention: postIntervention,
		})
	}
	return data, nil
}

// Calculate learning gains for control and intervention groups
func CalculateLearningGains(data []Data) (gains []float64, interventions []float64) {
	for _, d := range data {
		gainControl := d.PostControl - d.PreControl
		gainIntervention := d.PostIntervention - d.PreIntervention
		gains = append(gains, gainControl, gainIntervention)
		interventions = append(interventions, 0, 1) // 0 for control, 1 for intervention
	}
	return gains, interventions
}

// Linear regression function with Gonum
func LinearRegression(gains, interventions []float64) (beta0, beta1, rSquared float64) {
	weights := make([]float64, len(gains)) // Equal weights
	for i := range weights {
		weights[i] = 1.0
	}

	// Perform linear regression
	beta0, beta1 = stat.LinearRegression(interventions, gains, weights, false)
	rSquared = stat.RSquared(interventions, gains, weights, beta1, beta0)
	return beta0, beta1, rSquared
}

func ComputePValue(beta0, beta1 float64, gains, interventions []float64) float64 {
	// Calculate residuals using both intercept (beta0) and slope (beta1)
	residuals := make([]float64, len(gains))
	var sumSquaredResiduals float64
	for i := range gains {
		predicted := beta0 + beta1*interventions[i]
		residuals[i] = gains[i] - predicted
		sumSquaredResiduals += residuals[i] * residuals[i]
	}

	// Calculate the standard error of the slope (beta1)
	meanIntervention := stat.Mean(interventions, nil)
	sumSquaredDifferences := 0.0
	for _, x := range interventions {
		sumSquaredDifferences += (x - meanIntervention) * (x - meanIntervention)
	}
	standardError := math.Sqrt(sumSquaredResiduals / float64(len(gains)-2) / sumSquaredDifferences)

	// Calculate t-statistic
	tStatistic := beta1 / standardError

	// Print values for debugging
	log.Printf("[DEBUG] Residual Sum of Squares: %.8f\n", sumSquaredResiduals)
	log.Printf("[DEBUG] Standard Error: %.8f\n", standardError)
	log.Printf("[DEBUG] t-Statistic: %.8f\n", tStatistic)

	// Calculate p-value based on the t-distribution
	tDist := distuv.StudentsT{
		Mu:    0,
		Sigma: 1,
		Nu:    float64(len(gains) - 2),
	}
	pValue := 2 * (1 - tDist.CDF(math.Abs(tStatistic))) // Two-tailed test
	return pValue
}
