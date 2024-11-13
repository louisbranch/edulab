package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

func (srv *Server) participationsHandler(w http.ResponseWriter, r *http.Request,
	segments []string) {

	log.Print("[DEBUG] Routing participations")

	if len(segments) < 1 {
		srv.renderNotFound(w, r)
		return
	}

	pids := strings.Split(segments[0], "-")
	if len(pids) != 3 {
		srv.renderNotFound(w, r)
		return
	}

	eid, cid, aid := pids[0], pids[1], pids[2]

	experiment, err := srv.DB.FindExperiment(eid)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	cohort, err := srv.DB.FindCohort(experiment.ID, cid)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	assessment, err := srv.DB.FindAssessment(experiment.ID, aid)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	log.Printf("[DEBUG] Initiated participation for experiment %s, cohort %s, assessment %s",
		experiment.PublicID, cohort.PublicID, assessment.PublicID)

	var token string
	at, err := r.Cookie("access_token")
	if err == nil {
		token = at.Value
	} else {
		b := make([]byte, 8)
		srv.Random.Read(b)
		token = fmt.Sprintf("%x", b)
		http.SetCookie(w, &http.Cookie{
			Name:   "access_token",
			Value:  token,
			Path:   "/edulab/",
			MaxAge: 24 * 60 * 60 * 180, // 180 days
		})
	}

	participant, err := srv.DB.FindParticipant(experiment.ID, token)
	if err == sql.ErrNoRows {
		participant = edulab.Participant{
			PublicID:     srv.newPublicID(3),
			ExperimentID: experiment.ID,
			CohortID:     cohort.ID,
			AccessToken:  token,
		}
		err = srv.DB.CreateParticipant(&participant)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}
	}

	participations, err := srv.DB.FindParticipationsByParticipant(experiment.ID, participant.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	var participation edulab.Participation

	for _, p := range participations {
		if p.AssessmentID == assessment.ID {
			participation = p
			break
		}
		if len(p.Demographics) > 0 {
			participation.Demographics = p.Demographics
		}
	}

	if len(participation.Answers) > 0 {
		srv.participationCompleted(w, r, experiment, participant, assessment)
		return
	}

	if participation.Demographics == nil {
		demographics, err := srv.DB.FindDemographics(experiment.ID)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		if len(demographics) > 0 {
			srv.participateDemographics(w, r, experiment, participant, demographics)
			return
		}
	}

	srv.participateAssessment(w, r, experiment, participant, assessment)
}

func (srv *Server) participateAssessment(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, participant edulab.Participant,
	assessment edulab.Assessment) {

	printer, page := srv.i18n(w, r)

	questions, err := srv.DB.FindQuestions(assessment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	choices, err := srv.DB.FindQuestionChoices(assessment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	qp := presenter.GroupQuestions(questions, choices)

	page.Title = printer.Sprintf("%s - %s", experiment.Name, assessment.Type)
	page.Partials = []string{"assessment_participate"}
	page.Content = struct {
		Experiment  edulab.Experiment
		Participant edulab.Participant
		Assessment  presenter.Assessment
		Questions   []presenter.Question
		Texts       interface{}
	}{
		Experiment:  experiment,
		Participant: participant,
		Assessment:  presenter.NewAssessment(assessment, printer),
		Questions:   qp,
		Texts: struct {
			Submit string
		}{
			Submit: printer.Sprintf("Submit"),
		},
	}

	srv.render(w, page)
}

func (srv *Server) participateDemographics(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, participant edulab.Participant, demographics []edulab.Demographic) {

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
		Breadcrumbs  template.HTML
		Experiment   edulab.Experiment
		Participant  edulab.Participant
		Demographics []presenter.Demographic
		Texts        interface{}
	}{
		Breadcrumbs:  presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:   experiment,
		Participant:  participant,
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

func (srv *Server) participationCompleted(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, participant edulab.Participant,
	assessment edulab.Assessment) {

	printer, page := srv.i18n(w, r)

	page.Title = printer.Sprintf("Participation Completed")
	page.Partials = []string{"participation_completed"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Participant edulab.Participant
		Assessment  edulab.Assessment
		Texts       interface{}
	}{
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Participant: participant,
		Assessment:  assessment,
		Texts: struct {
			Title string
		}{
			Title: printer.Sprintf("Assessment Participation Completed"),
		},
	}

	srv.render(w, page)
}
