package wizard

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/louisbranch/edulab"
)

func Demographics(db edulab.Database, experiment edulab.Experiment) error {

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

	demographics_order := []string{gender, ageGroup, yearOfStudy, stemMajor}

	for _, category := range demographics_order {
		options := demographics[category]
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

	return nil
}
