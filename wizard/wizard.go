package wizard

import (
	"github.com/louisbranch/edulab"
)

type Experiment struct {
	PublicID        string          `yaml:"public_id"`
	Name            string          `yaml:"name"`
	Description     string          `yaml:"description"`
	Assessments     []Assessment    `yaml:"assessments"`
	Cohorts         []Cohort        `yaml:"cohorts"`
	BootstrapConfig BootstrapConfig `yaml:"bootstrap_config,omitempty"`
}

type Assessment struct {
	PublicID  string                `yaml:"public_id"`
	Type      edulab.AssessmentType `yaml:"type"`
	Questions []Question            `yaml:"questions"`
}

type Question struct {
	Text    string           `yaml:"text"`
	Type    edulab.InputType `yaml:"type"`
	Choices []Choice         `yaml:"choices"`
}

type Choice struct {
	Text      string `yaml:"text"`
	IsCorrect bool   `yaml:"is_correct"`
}

type Cohort struct {
	PublicID    string `yaml:"public_id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}
