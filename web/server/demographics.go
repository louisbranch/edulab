package server

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

// demographicsHandler handles the demographics subroutes in an experiment for the instructor.
func (srv *Server) demographicsHandler(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, segments []string) {

	log.Print("[DEBUG] web/server/demographics.go: handling demographics")

	if len(segments) < 1 {
		srv.listDemographics(w, r, experiment)
		return
	}

	switch segments[0] {
	case "preview":
		srv.previewDemographics(w, r, experiment)
		return
	default:
		srv.renderNotFound(w, r)
		return
	}
}

// listDemographics lists the demographics for an experiment to the instructor.
func (srv *Server) listDemographics(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	printer, page := srv.i18n(w, r)

	demographics, err := srv.DB.FindDemographics(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	title := printer.Sprintf("Demographics")
	page.Title = title
	page.Partials = []string{"demographics"}
	page.Content = struct {
		Breadcrumbs  template.HTML
		Experiment   edulab.Experiment
		Demographics []edulab.Demographic
		Texts        interface{}
	}{
		Breadcrumbs:  presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:   experiment,
		Demographics: demographics,
		Texts: struct {
			Title       string
			Demographic string
			Text        string
			Actions     string
			Add         string
			ComingSoon  string
			Preview     string
		}{
			Title:       title,
			Demographic: printer.Sprintf("Demographic"),
			Text:        printer.Sprintf("Text"),
			Actions:     printer.Sprintf("Actions"),
			Add:         printer.Sprintf("Add Demographic"),
			ComingSoon:  printer.Sprintf("Coming Soon"),
			Preview:     printer.Sprintf("Preview Questions"),
		},
	}

	srv.render(w, page)
}

// previewDemographics displays the demographics preview for the instructor.
func (srv *Server) previewDemographics(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	printer, page := srv.i18n(w, r)

	demographics, err := srv.DB.FindDemographics(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	dp := make(map[string]presenter.Demographic)
	for _, d := range demographics {
		dp[d.ID] = presenter.Demographic{
			Demographic: d,
		}
	}

	options, err := srv.DB.FindDemographicOptions(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	for _, o := range options {
		d, ok := dp[o.DemographicID]
		if !ok {
			log.Printf("[ERROR] web/server/assessments.go: demographic not found: %s", o.DemographicID)
			continue
		}
		d.Options = append(d.Options, o)
		dp[o.DemographicID] = d
	}

	page.Title = printer.Sprintf("Demographics")
	page.Partials = []string{"demographics_preview"}
	page.Content = struct {
		Breadcrumbs  template.HTML
		Experiment   edulab.Experiment
		Demographics []presenter.Demographic
		Texts        interface{}
	}{
		Breadcrumbs:  presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:   experiment,
		Demographics: presenter.SortDemographics(demographics, dp),
		Texts: struct {
			Title  string
			Submit string
			Back   string
		}{
			Title:  printer.Sprintf("Demographics"),
			Submit: printer.Sprintf("Submit"),
			Back:   printer.Sprintf("Back"),
		},
	}

	srv.render(w, page)
}

// showDemographics displays the demographics form for the participant.
func (srv *Server) showDemographics(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, cohort edulab.Cohort, participant edulab.Participant,
	assessment edulab.Assessment, demographics []edulab.Demographic) {

	printer, page := srv.i18n(w, r)

	dp := make(map[string]presenter.Demographic)
	for _, d := range demographics {
		dp[d.ID] = presenter.Demographic{
			Demographic: d,
		}
	}

	options, err := srv.DB.FindDemographicOptions(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	for _, o := range options {
		d, ok := dp[o.DemographicID]
		if !ok {
			log.Printf("[ERROR] web/server/assessments.go: demographic not found: %s", o.DemographicID)
			continue
		}
		d.Options = append(d.Options, o)
		dp[o.DemographicID] = d
	}

	page.Title = printer.Sprintf("Demographics")
	page.Partials = []string{"demographics_participate"}
	page.Content = struct {
		Breadcrumbs template.HTML
		edulab.Experiment
		edulab.Participant
		edulab.Assessment
		edulab.Cohort
		Demographics []presenter.Demographic
		Texts        interface{}
	}{
		Breadcrumbs:  presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:   experiment,
		Assessment:   assessment,
		Participant:  participant,
		Cohort:       cohort,
		Demographics: presenter.SortDemographics(demographics, dp),
		Texts: struct {
			Title  string
			Submit string
		}{
			Title:  printer.Sprintf("Demographics"),
			Submit: printer.Sprintf("Submit"),
		},
	}

	srv.render(w, page)
}

func (srv *Server) participateDemographics(w http.ResponseWriter, r *http.Request) {
	log.Print("[DEBUG] Request to participate in demographics")

	if r.Method != "POST" {
		srv.renderNotFound(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	form := r.PostForm

	eid := form.Get("experiment_id")
	aid := form.Get("assessment_id")
	cid := form.Get("cohort_id")
	token := form.Get("participant_access_token")

	inputs := make(map[string][]string)

	for key, values := range form {
		found := false

		for _, ignore := range []string{"experiment_id", "assessment_id", "cohort_id", "participant_access_token"} {
			if key == ignore {
				found = true
				break
			}
		}

		if !found {
			inputs[key] = values
		}
	}

	demographics, err := marshalToSortedRawMessage(inputs)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	experiment, err := srv.DB.FindExperiment(eid)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	assessment, err := srv.DB.FindAssessment(experiment.ID, aid)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	participant, err := srv.DB.FindParticipant(experiment.ID, token)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	participation, err := srv.DB.FindParticipation(experiment.ID, assessment.ID, participant.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		srv.renderError(w, r, err)
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		participation = edulab.Participation{
			ExperimentID:  experiment.ID,
			AssessmentID:  assessment.ID,
			ParticipantID: participant.ID,
			Demographics:  demographics,
		}

		err = srv.DB.CreateParticipation(&participation)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}
	} else {
		participation.Demographics = demographics

		err = srv.DB.UpdateParticipation(participation)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/edulab/%s-%s-%s", eid, cid, aid), http.StatusTemporaryRedirect)
}
