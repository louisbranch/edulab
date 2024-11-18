package wizard

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/louisbranch/edulab"
)

// ImportYAML loads and imports all YAML experiment files from a directory.
func ImportYAML(db edulab.Database, dirPath string) error {
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "error accessing path %s", path)
		}

		// Process only files with .yaml extension
		if d.Type().IsRegular() && filepath.Ext(path) == ".yaml" {
			log.Printf("[INFO] Importing experiment from file: %s\n", path)

			experiment, err := loadYAML(path)
			if err != nil {
				return errors.Wrapf(err, "error loading experiment from %s", path)
			}

			if err := create(db, experiment); err != nil {
				return errors.Wrapf(err, "error importing experiment from %s", path)
			}
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "error walking directory")
	}

	log.Println("[INFO] All experiments have been imported successfully.")
	return nil
}

func loadYAML(filename string) (Experiment, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Experiment{}, errors.Wrap(err, "could not open YAML file")
	}
	defer file.Close()

	var experiment Experiment
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&experiment); err != nil {
		return Experiment{}, errors.Wrap(err, "could not decode YAML file")
	}

	return experiment, nil
}

func create(db edulab.Database, experimentData Experiment) error {

	// Check if experiment already exists
	experiment, err := db.FindExperiment(experimentData.PublicID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(err, "could not find experiment")
	}

	if experimentData.PublicID != "" && experimentData.ForceDelete {
		log.Printf("[INFO] Experiment %s already exists, deleting and recreating.\n", experimentData.PublicID)
		err = db.DeleteExperiment(experiment.ID)
		if err != nil {
			return errors.Wrap(err, "could not delete experiment")
		}
	} else if experiment.PublicID != "" && !experimentData.ForceDelete {
		log.Printf("[INFO] Experiment %s already exists, skipping creation.\n", experimentData.PublicID)
		err = bootstrapParticipants(db, experimentData.BootstrapConfig, experiment)
		if err != nil {
			return errors.Wrap(err, "could not bootstrap participants")
		}
		return nil
	}

	experiment = edulab.Experiment{
		PublicID:    experimentData.PublicID,
		Name:        experimentData.Name,
		Description: experimentData.Description,
		CreatedAt:   time.Now(),
	}

	if err := db.CreateExperiment(&experiment); err != nil {
		return errors.Wrap(err, "could not create experiment")
	}

	for _, a := range experimentData.Assessments {
		assessment := edulab.Assessment{
			PublicID:     a.PublicID,
			ExperimentID: experiment.ID,
			Type:         a.Type,
		}
		if err := db.CreateAssessment(&assessment); err != nil {
			return errors.Wrap(err, "could not create assessment")
		}

		for _, q := range a.Questions {
			question := edulab.Question{
				AssessmentID: assessment.ID,
				Text:         q.Text,
				Type:         q.Type,
			}
			if err := db.CreateQuestion(&question); err != nil {
				return errors.Wrap(err, "could not create question")
			}

			for _, choice := range q.Choices {
				questionChoice := edulab.QuestionChoice{
					QuestionID: question.ID,
					Text:       choice.Text,
					IsCorrect:  choice.IsCorrect,
				}
				if err := db.CreateQuestionChoice(&questionChoice); err != nil {
					return errors.Wrap(err, "could not create question choice")
				}
			}
		}
	}

	// Import cohorts
	for _, c := range experimentData.Cohorts {
		cohort := edulab.Cohort{
			PublicID:     c.PublicID,
			ExperimentID: experiment.ID,
			Name:         c.Name,
			Description:  c.Description,
		}

		if err := db.CreateCohort(&cohort); err != nil {
			return errors.Wrap(err, "could not create cohort")
		}
	}

	// Import demographics
	err = Demographics(db, experiment)
	if err != nil {
		return errors.Wrap(err, "could not create demographics")
	}

	// Import bootstrap configuration
	if experimentData.BootstrapConfig.Participants == 0 {
		return nil
	}

	err = bootstrapParticipants(db, experimentData.BootstrapConfig, experiment)
	if err != nil {
		return errors.Wrap(err, "could not bootstrap participants")
	}

	return nil
}
