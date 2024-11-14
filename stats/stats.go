package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

// Data struct to hold the pre- and post-intervention scores
type Data struct {
	PreBase          float64
	PostBase         float64
	PreIntervention  float64
	PostIntervention float64
}

// Function to read data from CSV file
func readCSV(filePath string) ([]Data, error) {
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
		preBase, _ := strconv.ParseFloat(record[0], 64)
		preIntervention, _ := strconv.ParseFloat(record[1], 64)
		postBase, _ := strconv.ParseFloat(record[2], 64)
		postIntervention, _ := strconv.ParseFloat(record[3], 64)

		data = append(data, Data{
			PreBase:          preBase,
			PostBase:         postBase,
			PreIntervention:  preIntervention,
			PostIntervention: postIntervention,
		})
	}
	return data, nil
}

// Calculate learning gains for base and intervention groups
func calculateLearningGains(data []Data) (gains []float64, interventions []float64) {
	for _, d := range data {
		gainBase := d.PostBase - d.PreBase
		gainIntervention := d.PostIntervention - d.PreIntervention
		gains = append(gains, gainBase, gainIntervention)
		interventions = append(interventions, 0, 1) // 0 for base, 1 for intervention
	}
	return gains, interventions
}

// Linear regression function with Gonum
func linearRegression(gains, interventions []float64) (beta0, beta1, rSquared float64) {
	weights := make([]float64, len(gains)) // Equal weights
	for i := range weights {
		weights[i] = 1.0
	}

	// Perform linear regression
	beta0, beta1 = stat.LinearRegression(interventions, gains, weights, false)
	rSquared = stat.RSquared(interventions, gains, weights, beta1, beta0)
	return beta0, beta1, rSquared
}

func computePValue(beta0, beta1 float64, gains, interventions []float64) float64 {
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
	fmt.Printf("Residual Sum of Squares: %.5f\n", sumSquaredResiduals)
	fmt.Printf("Standard Error: %.5f\n", standardError)
	fmt.Printf("t-Statistic: %.5f\n", tStatistic)

	// Calculate p-value based on the t-distribution
	tDist := distuv.StudentsT{
		Mu:    0,
		Sigma: 1,
		Nu:    float64(len(gains) - 2),
	}
	pValue := 2 * (1 - tDist.CDF(math.Abs(tStatistic))) // Two-tailed test
	return pValue
}

func main() {
	// Use flag to accept the CSV file path as a command-line argument
	filePath := flag.String("file", "", "Path to the CSV file containing data")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Please provide a file path using the -file flag.")
		return
	}

	// Load the CSV data
	data, err := readCSV(*filePath)
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Calculate learning gains
	gains, interventions := calculateLearningGains(data)

	// Perform linear regression
	beta0, beta1, rSquared := linearRegression(gains, interventions)
	fmt.Printf("Intercept (beta0): %.5f, Slope (beta1): %.5f, R-squared: %.5f\n", beta0, beta1, rSquared)

	// Calculate p-value for the slope
	pValue := computePValue(beta0, beta1, gains, interventions)
	fmt.Printf("P-value for the intervention coefficient: %.5f\n", pValue)

	// Interpret the p-value
	alpha := 0.2
	if pValue < alpha {
		fmt.Println("The intervention effect is statistically significant (p < alpha).")
	} else {
		fmt.Println("The intervention effect is not statistically significant (p >= alpha).")
	}
}
