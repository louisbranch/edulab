package main

import (
	"flag"
	"fmt"

	"github.com/louisbranch/edulab/stats"
)

func main() {
	// Use flag to accept the CSV file path as a command-line argument
	filePath := flag.String("file", "", "Path to the CSV file containing data")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Please provide a file path using the -file flag.")
		return
	}

	// Load the CSV data
	data, err := stats.ReadCSV(*filePath)
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Calculate learning gains
	gains, interventions := stats.CalculateLearningGains(data)

	// Perform linear regression
	beta0, beta1, rSquared := stats.LinearRegression(gains, interventions)
	fmt.Printf("Intercept (beta0): %.8f, Slope (beta1): %.8f, R-squared: %.8f\n", beta0, beta1, rSquared)

	// Calculate p-value for the slope
	pValue := stats.ComputePValue(beta0, beta1, gains, interventions)
	fmt.Printf("P-value for the intervention coefficient: %.8f\n", pValue)

	// Interpret the p-value
	alpha := 0.2
	if pValue < alpha {
		fmt.Println("The intervention effect is statistically significant (p < alpha).")
	} else {
		fmt.Println("The intervention effect is not statistically significant (p >= alpha).")
	}
}
