package wizard

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/louisbranch/edulab"
)

func Experiment(db edulab.Database) error {

	experiments, err := db.FindExperiments()
	if err != nil {
		return errors.Wrap(err, "wizard could not find experiments")
	}

	if len(experiments) > 0 {
		return nil
	}

	printer := message.NewPrinter(language.English)

	experiment := edulab.Experiment{
		PublicID: "NCC-1701",
		Name:     "Earth's Seasons",
		Description: `This experiment gauges students' understanding of the Earth's seasons.
One cohort will attend a traditional lecture, while the other will attend a workshop.`,
	}

	err = db.CreateExperiment(&experiment)
	if err != nil {
		return errors.Wrap(err, "wizard could not create experiment")
	}

	assessments := []edulab.Assessment{}
	assessmentsTypes := []edulab.AssessmentType{edulab.AssessmentTypePre, edulab.AssessmentTypePos}
	for i, assessmentType := range assessmentsTypes {
		assessment := edulab.Assessment{
			PublicID:     fmt.Sprintf("A-%d", i),
			ExperimentID: experiment.ID,
			Type:         assessmentType,
		}
		err = db.CreateAssessment(&assessment)
		if err != nil {
			return errors.Wrap(err, "wizard could not create assessment")
		}

		assessments = append(assessments, assessment)
	}

	type questionChoice struct {
		text       string
		is_correct bool
	}

	type question struct {
		text    string
		qtype   string // single, multiple, text
		choices []questionChoice
	}

	questions := []question{
		{
			text:  "What is the primary cause of the **seasons** on Earth?",
			qtype: "single",
			choices: []questionChoice{
				{text: "The **distance** between the Earth and the Sun changes throughout the year."},
				{text: "The Earth's **axis** is **tilted** as it **orbits** the Sun.", is_correct: true},
				{text: "The **speed** of Earth's orbit changes throughout the year."},
				{text: "Different parts of Earth receive **sunlight** based on its **rotation**."},
			},
		},
		{
			text:  "Which of the following statements are **true** about Earth’s **orbit** and **seasons**? _(Select all that apply)_",
			qtype: "multiple",
			choices: []questionChoice{
				{text: "Earth is **closer** to the Sun in summer."},
				{text: "Earth’s **tilt** causes different parts of Earth to receive varying amounts of **sunlight**.", is_correct: true},
				{text: "Earth’s **distance** from the Sun varies greatly, causing seasons."},
				{text: "The **tilt** of Earth's axis remains **constant** relative to its orbit.", is_correct: true},
			},
		},
		{
			text:  "In your own words, describe why it is **warmer** in **summer** than in **winter**. *(Maximum 100 characters)*",
			qtype: "text",
		},
		{
			text:  "When it is **summer** in the **Northern Hemisphere**, what season is it in the **Southern Hemisphere**?",
			qtype: "single",
			choices: []questionChoice{
				{text: "Winter", is_correct: true},
				{text: "Spring"},
				{text: "Summer"},
				{text: "Autumn"},
			},
		},
		{
			text:  "How does the **tilt** of the Earth affect the **intensity** and **duration** of **sunlight** received at different locations on Earth? _(Select all that apply)_",
			qtype: "multiple",
			choices: []questionChoice{
				{text: "The **tilt** changes which hemisphere is tilted towards the **Sun**, affecting sunlight **intensity**.", is_correct: true},
				{text: "The **tilt** causes Earth to be **closer** to the Sun during different seasons."},
				{text: "The **tilt** affects the **angle** and **duration** of sunlight, influencing **temperatures** and **day length**.", is_correct: true},
			},
		},
	}

	for _, assessment := range assessments {
		for _, q := range questions {

			qs := edulab.Question{
				AssessmentID: assessment.ID,
				Prompt:       q.text,
				Type:         edulab.InputType(q.qtype),
			}

			err := db.CreateQuestion(&qs)
			if err != nil {
				return errors.Wrap(err, "wizard could not create question")
			}

			for _, c := range q.choices {
				err := db.CreateQuestionChoice(&edulab.QuestionChoice{
					QuestionID: qs.ID,
					Text:       c.text,
					IsCorrect:  c.is_correct,
				})
				if err != nil {
					return errors.Wrap(err, "wizard could not create question choice")
				}
			}
		}
	}

	type i18npair struct {
		key string
		val string
	}

	demographics := map[i18npair][]i18npair{
		{key: "gender", val: printer.Sprint("Gender")}: {
			{key: "male", val: printer.Sprint("Male")},
			{key: "female", val: printer.Sprint("Female")},
			{key: "non_binary", val: printer.Sprint("Non-binary")},
			{key: "prefer_not_to_say", val: printer.Sprint("Prefer not to say")},
		},

		{key: "age_group", val: printer.Sprint("Age Group")}: {
			{key: "<18", val: printer.Sprint("Under 18")},
			{key: "18_20", val: printer.Sprint("18 to 20")},
			{key: "21_23", val: printer.Sprint("21 to 23")},
			{key: "24_26", val: printer.Sprint("24 to 26")},
			{key: "27+", val: printer.Sprint("27+")},
		},

		{key: "year_study", val: printer.Sprint("Year of Study")}: {
			{key: "1", val: printer.Sprint("Year 1")},
			{key: "2", val: printer.Sprint("Year 2")},
			{key: "3", val: printer.Sprint("Year 3")},
			{key: "4", val: printer.Sprint("Year 4")},
			{key: "5+", val: printer.Sprint("Year 5+")},
		},

		{key: "stem_major", val: printer.Sprint("STEM Major")}: {
			{key: "physical_sciences", val: printer.Sprint("Physical Sciences")},
			{key: "life_sciences", val: printer.Sprint("Life Sciences")},
			{key: "earth_environmental_sciences", val: printer.Sprint("Earth & Environmental Sciences")},
			{key: "mathematics_computer_science", val: printer.Sprint("Mathematics & Computer Science")},
			{key: "engineering", val: printer.Sprint("Engineering")},
			{key: "other", val: printer.Sprint("Other")},
		},
	}

	for category, options := range demographics {
		d := edulab.Demographic{
			ExperimentID: experiment.ID,
			I18nKey:      category.key,
			Type:         edulab.InputSingle,
		}

		err := db.CreateDemographic(&d)
		if err != nil {
			return err
		}

		for _, option := range options {
			err := db.CreateDemographicOption(&edulab.DemographicOption{
				DemographicID: d.ID,
				I18nKey:       option.key,
			})
			if err != nil {
				return err
			}
		}
	}

	cohorts := []string{"Baseline", "Workshop"}
	for i, cohort := range cohorts {
		err = db.CreateCohort(&edulab.Cohort{
			PublicID:     fmt.Sprintf("C-%d", i),
			ExperimentID: experiment.ID,
			Name:         cohort,
		})
		if err != nil {
			return errors.Wrap(err, "wizard could not create cohort")
		}
	}

	return nil
}
