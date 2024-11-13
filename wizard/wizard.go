package wizard

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/louisbranch/edulab"
)

func Experiment(db edulab.Database) error {

	printer := message.NewPrinter(language.English)

	var (
		gender       = printer.Sprintf("Gender")
		male         = printer.Sprintf("Male")
		female       = printer.Sprintf("Female")
		nonBinary    = printer.Sprintf("Non-binary")
		preferNotSay = printer.Sprintf("Prefer not to say")

		ageGroup  = printer.Sprintf("Age Group")
		under18   = printer.Sprintf("Under 18")
		age18To20 = printer.Sprintf("18 to 20")
		age21To23 = printer.Sprintf("21 to 23")
		age24To26 = printer.Sprintf("24 to 26")
		age27Plus = printer.Sprintf("27+")

		yearOfStudy = printer.Sprintf("Year of Study")
		year1       = printer.Sprintf("Year 1")
		year2       = printer.Sprintf("Year 2")
		year3       = printer.Sprintf("Year 3")
		year4       = printer.Sprintf("Year 4")
		year5Plus   = printer.Sprintf("Year 5+")

		stemMajor                  = printer.Sprintf("STEM Major")
		physicalSciences           = printer.Sprintf("Physical Sciences")
		lifeSciences               = printer.Sprintf("Life Sciences")
		earthEnvironmentalSciences = printer.Sprintf("Earth & Environmental Sciences")
		mathematicsComputerScience = printer.Sprintf("Mathematics & Computer Science")
		engineering                = printer.Sprintf("Engineering")
		other                      = printer.Sprintf("Other")
	)

	var demographics = map[string][]string{
		gender: {
			male,
			female,
			nonBinary,
			preferNotSay,
		},

		ageGroup: {
			under18,
			age18To20,
			age21To23,
			age24To26,
			age27Plus,
		},

		yearOfStudy: {
			year1,
			year2,
			year3,
			year4,
			year5Plus,
		},

		stemMajor: {
			physicalSciences,
			lifeSciences,
			earthEnvironmentalSciences,
			mathematicsComputerScience,
			engineering,
			other,
		},
	}

	experiments, err := db.FindExperiments()
	if err != nil {
		return errors.Wrap(err, "wizard could not find experiments")
	}

	if len(experiments) > 0 {
		return nil
	}

	experiment := edulab.Experiment{
		PublicID: "E1",
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
			PublicID:     fmt.Sprintf("A%d", i+1),
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
				Text:         q.text,
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

	for category, options := range demographics {
		d := edulab.Demographic{
			ExperimentID: experiment.ID,
			Text:         category,
			Type:         edulab.InputSingle,
		}

		err := db.CreateDemographic(&d)
		if err != nil {
			return err
		}

		for _, option := range options {
			err := db.CreateDemographicOption(&edulab.DemographicOption{
				DemographicID: d.ID,
				Text:          option,
			})
			if err != nil {
				return err
			}
		}
	}

	cohorts := []string{"Baseline", "Workshop"}
	for i, cohort := range cohorts {
		err = db.CreateCohort(&edulab.Cohort{
			PublicID:     fmt.Sprintf("C%d", i+1),
			ExperimentID: experiment.ID,
			Name:         cohort,
		})
		if err != nil {
			return errors.Wrap(err, "wizard could not create cohort")
		}
	}

	return nil
}
