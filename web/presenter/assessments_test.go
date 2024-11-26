package presenter

import (
	"reflect"
	"testing"

	"github.com/louisbranch/edulab"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestNewAssessment(t *testing.T) {
	printer := message.NewPrinter(language.English)

	tests := []struct {
		name     string
		input    edulab.Assessment
		expected Assessment
	}{
		{
			name: "Pre-Assessment",
			input: edulab.Assessment{
				Type: edulab.AssessmentTypePre,
			},
			expected: Assessment{
				Assessment: edulab.Assessment{
					Type: edulab.AssessmentTypePre,
				},
				Type: "Pre-Assessment",
			},
		},
		{
			name: "Post-Assessment",
			input: edulab.Assessment{
				Type: edulab.AssessmentTypePost,
			},
			expected: Assessment{
				Assessment: edulab.Assessment{
					Type: edulab.AssessmentTypePost,
				},
				Type: "Post-Assessment",
			},
		},
		{
			name: "Unknown Assessment Type",
			input: edulab.Assessment{
				Type: edulab.AssessmentType("Unknown"),
			},
			expected: Assessment{
				Assessment: edulab.Assessment{
					Type: edulab.AssessmentType("Unknown"),
				},
				Type: "Unknown Assessment Type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewAssessment(tt.input, printer)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("NewAssessment() = %v, want %v", actual, tt.expected)
			}
		})
	}
}
func TestNewAssessments(t *testing.T) {
	printer := message.NewPrinter(language.English)

	tests := []struct {
		name     string
		input    []edulab.Assessment
		expected []Assessment
	}{
		{
			name: "Multiple Assessments",
			input: []edulab.Assessment{
				{Type: edulab.AssessmentTypePre},
				{Type: edulab.AssessmentTypePost},
				{Type: edulab.AssessmentType("Unknown")},
			},
			expected: []Assessment{
				{
					Assessment: edulab.Assessment{Type: edulab.AssessmentTypePre},
					Type:       "Pre-Assessment",
				},
				{
					Assessment: edulab.Assessment{Type: edulab.AssessmentTypePost},
					Type:       "Post-Assessment",
				},
				{
					Assessment: edulab.Assessment{Type: edulab.AssessmentType("Unknown")},
					Type:       "Unknown Assessment Type",
				},
			},
		},
		{
			name:     "Empty Assessments",
			input:    []edulab.Assessment{},
			expected: []Assessment{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewAssessments(tt.input, printer)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("NewAssessments() = %v, want %v", actual, tt.expected)
			}
		})
	}
}
func TestAssessmentType(t *testing.T) {
	printer := message.NewPrinter(language.English)

	tests := []struct {
		name     string
		input    edulab.AssessmentType
		expected string
	}{
		{
			name:     "Pre-Assessment",
			input:    edulab.AssessmentTypePre,
			expected: "Pre-Assessment",
		},
		{
			name:     "Post-Assessment",
			input:    edulab.AssessmentTypePost,
			expected: "Post-Assessment",
		},
		{
			name:     "Unknown Assessment Type",
			input:    edulab.AssessmentType("Unknown"),
			expected: "Unknown Assessment Type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := AssessmentType(printer, tt.input)
			if actual != tt.expected {
				t.Errorf("AssessmentType() = %v, want %v", actual, tt.expected)
			}
		})
	}
}
func TestAssessmentTypes(t *testing.T) {
	printer := message.NewPrinter(language.English)

	expected := [][]string{
		{string(edulab.AssessmentTypePre), "Pre-Assessment"},
		{string(edulab.AssessmentTypePost), "Post-Assessment"},
	}

	actual := AssessmentTypes(printer)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("AssessmentTypes() = %v, want %v", actual, expected)
	}
}
