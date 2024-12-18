package server

import (
	"encoding/json"
	"html/template"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/louisbranch/edulab"
	"github.com/louisbranch/edulab/result"
	"github.com/louisbranch/edulab/stats"
	"github.com/louisbranch/edulab/web/presenter"
	"gonum.org/v1/gonum/stat"
)

// Naive cache for learning gains results
// TODO: Implement a proper cache
type cached struct {
	experimentID   string
	participations int
	payload        []byte
}

var gainsCache = make(map[string]cached)

func (srv *Server) resultsHandler(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment, segments []string) {

	log.Print("[DEBUG] Routing results")

	if len(segments) < 1 {
		srv.renderNotFound(w, r)
		return
	}

	switch segments[0] {
	case "demographics":
		srv.demographicsResult(w, r, experiment)
		return
	case "assessments":
		srv.assessmentsResult(w, r, experiment)
		return
	case "gains":
		srv.gainsResult(w, r, experiment)
		return
	default:
		srv.renderNotFound(w, r)
		return
	}
}

func (srv *Server) demographicsResult(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	cohorts, err := srv.DB.FindCohorts(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	demographics, err := srv.DB.FindDemographics(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	options, err := srv.DB.FindDemographicOptions(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	participants, err := srv.DB.FindParticipants(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	participations, err := srv.DB.FindParticipations(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	dr, err := presenter.NewDemographicsResult(demographics, options, cohorts,
		participants, participations)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	if r.Header.Get("Content-type") == "application/json" {

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(dr.Data)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		return
	}

	printer, page := srv.i18n(w, r)

	title := printer.Sprintf("Demographics Results")
	page.Title = title
	page.Partials = []string{"results_demographics"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Results     presenter.DemographicsResult
		Texts       interface{}
	}{
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Results:     dr,
		Texts: struct {
			Title        string
			Download     string
			ComingSoon   string
			Options      string
			Participants string
			Empty        string
		}{
			Title:        title,
			Download:     printer.Sprintf("Export as CSV"),
			ComingSoon:   printer.Sprintf("Coming soon"),
			Options:      printer.Sprintf("Options"),
			Participants: printer.Sprintf("Participants"),
			Empty:        printer.Sprintf("No data available yet"),
		},
	}

	srv.render(w, page)
}

func (srv *Server) assessmentsResult(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	printer, page := srv.i18n(w, r)

	assessments, err := srv.DB.FindAssessments(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	participations, err := srv.DB.FindParticipations(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	var allChoices [][]edulab.QuestionChoice

	aps := presenter.NewAssessments(assessments, printer)
	for i := range aps {
		questions, err := srv.DB.FindQuestions(assessments[i].ID)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		choices, err := srv.DB.FindQuestionChoices(assessments[i].ID)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		for _, q := range questions {
			if q.Type == edulab.InputText {
				continue
			}

			var cs []edulab.QuestionChoice
			for _, c := range choices {
				if c.QuestionID != q.ID {
					continue
				}
				cs = append(cs, c)
			}
			allChoices = append(allChoices, cs)
		}

		aps[i].Questions = presenter.GroupQuestions(questions, choices)

		for _, p := range participations {
			if p.AssessmentID != assessments[i].ID {
				continue
			}

		}
	}

	if r.Header.Get("Content-type") == "application/json" {

		counts, err := result.CountChoicesByCohorts(srv.DB, experiment)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(counts)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		return
	}

	title := printer.Sprintf("Assessments Results")
	page.Title = title
	page.Partials = []string{"results_assessments"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Assessments []presenter.Assessment
		Choices     [][]edulab.QuestionChoice
		Texts       interface{}
	}{
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Assessments: aps,
		Choices:     allChoices,
		Texts: struct {
			Title        string
			Choices      string
			Participants string
			Empty        string
			CohortLabels []string
		}{
			Title:        title,
			Choices:      printer.Sprintf("Choices"),
			Participants: printer.Sprintf("Participants"),
			Empty:        printer.Sprintf("No data available yet"),
			CohortLabels: []string{
				printer.Sprintf("Control"),
				printer.Sprintf("Intervention"),
			},
		},
	}

	srv.render(w, page)
}

func (srv *Server) gainsResult(w http.ResponseWriter, r *http.Request,
	experiment edulab.Experiment) {

	res, err := result.New(srv.DB, experiment.ID)
	if err != nil {
		log.Printf("[ERROR] Failed to create result: %v", err)
		srv.renderError(w, r, err)
		return
	}

	printer, page := srv.i18n(w, r)

	type texts struct {
		Title           string
		Error           string
		Download        string
		ComingSoon      string
		Empty           string
		PlotTitles      []string
		AssessmentTypes []string
		CohortLabels    []string
	}

	content := struct {
		Breadcrumbs template.HTML
		Experiment  edulab.Experiment
		Texts       texts
	}{
		Breadcrumbs: presenter.ExperimentBreadcrumb(experiment, printer),
		Experiment:  experiment,
		Texts: texts{
			Title:      printer.Sprintf("Gains Results"),
			Download:   printer.Sprintf("Export as CSV"),
			ComingSoon: printer.Sprintf("Coming soon"),
			PlotTitles: []string{
				printer.Sprintf("Average Correct Answers by Cohort"),
				printer.Sprintf("Learning Gain by Cohort (Post - Pre)"),
			},
			AssessmentTypes: []string{
				printer.Sprintf("Pre"),
				printer.Sprintf("Post"),
			},
			CohortLabels: []string{
				printer.Sprintf("Control"),
				printer.Sprintf("Intervention"),
			},
		},
	}

	title := printer.Sprintf("Gains Results")
	page.Title = title
	page.Partials = []string{"results_gains"}
	page.Content = content

	cache, ok := gainsCache[experiment.ID]
	if ok && cache.participations == res.Participations() {
		if r.Header.Get("Content-type") != "application/json" {
			srv.render(w, page)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(cache.payload)
		return
	}

	if !res.Valid() {
		content.Texts.Error = printer.Sprintf("No data available yet")
		page.Content = content
		srv.render(w, page)
		return
	}

	err = res.Load()
	if err != nil {
		log.Printf("[ERROR] Failed to load result: %v", err)
		srv.renderError(w, r, err)
		return
	}

	cohorts, items := res.ComparisonPairs()
	if len(items) == 0 {
		content.Texts.Error = printer.Sprintf("No comparison pairs available yet")
		page.Content = content
		srv.render(w, page)
		return
	}

	if r.Header.Get("Content-type") != "application/json" {
		srv.render(w, page)
		return
	}

	participants, err := srv.DB.FindParticipants(experiment.ID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	questions, err := srv.DB.FindQuestions(items[0][0].AssessmentID)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	labels := make([]string, len(questions))
	for i, q := range questions {
		s := q.Text

		for _, char := range []rune{'\n', '\r', '\t', '*', '_'} {
			s = strings.ReplaceAll(s, string(char), "")
		}

		if len(s) <= 200 {
			labels[i] = s
			continue
		}
		labels[i] = s[:200] + "..."
	}

	type chart struct {
		Question         string  `json:"question"`
		PreControl       float64 `json:"preControl"`
		PostControl      float64 `json:"postControl"`
		PreIntervention  float64 `json:"preIntervention"`
		PostIntervention float64 `json:"postIntervention"`
		Beta0            float64 `json:"beta0"`
		Beta1            float64 `json:"beta1"`
		RSquared         float64 `json:"rSquared"`
		PValue           float64 `json:"pValue"`
		Message          string  `json:"message"`
	}

	var payload []chart

	for i, item := range items {

		var label string
		if len(labels) > i {
			label = labels[i]
		}

		comparison, err := result.NewComparison(res, item, cohorts)
		if err != nil {
			srv.renderError(w, r, err)
			return
		}

		data := comparison.ToStatsData()

		gains, interventions := stats.CalculateLearningGains(data)

		beta0, beta1, rSquared := stats.LinearRegression(gains, interventions)

		pValue := stats.ComputePValue(beta0, beta1, gains, interventions)

		if math.IsNaN(pValue) {
			pValue = 1.0
		}

		if math.IsNaN(rSquared) {
			rSquared = 0.0
		}

		var preControl, postControl, preIntervention, postIntervention []float64
		for _, d := range data {
			preControl = append(preControl, d.PreControl)
			postControl = append(postControl, d.PostControl)
			preIntervention = append(preIntervention, d.PreIntervention)
			postIntervention = append(postIntervention, d.PostIntervention)
		}

		payload = append(payload, chart{
			Question:         label,
			PreControl:       stat.Mean(preControl, nil),
			PostControl:      stat.Mean(postControl, nil),
			PreIntervention:  stat.Mean(preIntervention, nil),
			PostIntervention: stat.Mean(postIntervention, nil),
			Beta0:            beta0,
			Beta1:            beta1,
			RSquared:         rSquared,
			PValue:           pValue,
			Message:          result.EvaluateExperiment(len(participants), pValue, printer),
		})

	}

	w.Header().Set("Content-Type", "application/json")
	// Encode the payload into a []byte
	response, err := json.Marshal(payload)
	if err != nil {
		srv.renderError(w, r, err)
		return
	}

	gainsCache[experiment.ID] = cached{
		experimentID:   experiment.ID,
		participations: res.Participations(),
		payload:        response,
	}

	w.Write(response)
}
