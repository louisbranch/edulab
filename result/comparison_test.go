package result

import (
	"os"
	"testing"
)

func TestToCSV(t *testing.T) {

	c := &Comparison{
		headers: []string{"header1", "header2"},
		data: map[string][]float64{
			"header1": {1.0, 2.0},
			"header2": {3.0, 4.0},
		},
		rows: 2,
	}

	tmpFile, err := os.CreateTemp("", "comparison-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file after the test

	err = c.ToCSV(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to write CSV file: %v", err)
	}

	// Read the file
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open temporary file: %v", err)
	}
	defer file.Close()

	csvData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temporary file: %v", err)
	}

	// Check the output
	expected := "header1,header2\n1.00,3.00\n2.00,4.00\n"
	if string(csvData) != expected {
		t.Errorf("Expected %q, got %q", expected, csvData)
	}

}
