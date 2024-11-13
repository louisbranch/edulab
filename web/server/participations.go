package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/web/presenter"
)

// participationsHandler handles the participations subroutes in an experiment for the participant.
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
	if errors.Is(err, sql.ErrNoRows) {
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
		srv.participationFinished(w, r, experiment, participant, assessment)
		return
	}

	if participation.Demographics == nil {
		demographics, err := srv.DB.FindDemographics(experiment.ID)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		if len(demographics) > 0 {
			srv.showDemographics(w, r, experiment, cohort, participant, assessment, demographics)
			return
		}
	}

	srv.showAssessment(w, r, experiment, cohort, participant, assessment)
}

// participationFinished displays the participation finished page.
func (srv *Server) participationFinished(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, participant edulab.Participant,
	assessment edulab.Assessment) {

	printer, page := srv.i18n(w, r)

	title := printer.Sprintf("Thank you for participating!")
	page.Title = title
	page.Partials = []string{"participation_finished"}
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
			Title   string
			Message string
		}{
			Title:   title,
			Message: printer.Sprintf("Your participation has been successfully recorded.\n\nYou can now close this page."),
		},
	}

	srv.render(w, page)
}

// marshalToSortedRawMessage marshals a map to a JSON raw message with sorted keys.
// Used to serialize participation data for storage.
func marshalToSortedRawMessage(data map[string][]string) (json.RawMessage, error) {
	// Extract and sort the keys
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Create a temporary ordered map
	orderedData := make(map[string][]string, len(data))
	for _, key := range keys {
		orderedData[key] = data[key]
	}

	// Marshal the ordered map to JSON and return as json.RawMessage
	jsonData, err := json.Marshal(orderedData)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(jsonData), nil
}
