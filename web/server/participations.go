package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/louisbranch/edulab"
)

func (srv *Server) participationsHandler(w http.ResponseWriter, r *http.Request,
	segments []string) {

	log.Print("[DEBUG] web/server/participations.go: handling participations")

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

	fmt.Fprintf(w, "Participation: %s, %s, %s", experiment.Name, cohort.Name, assessment.Type)
}
