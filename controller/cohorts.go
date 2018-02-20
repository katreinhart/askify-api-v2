package controller

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/katreinhart/askify-api-v2/model"
)

// FetchCohortList responds a JSON object of a list of cohorts.
func FetchCohortList(w http.ResponseWriter, r *http.Request) {
	var cohorts []model.Cohort

	cohorts = model.FetchCohortList()
	js, err := json.Marshal(cohorts)
	handleErrorAndRespond(js, err, w)
}

// AddCohort adds a new cohort to the database list
func AddCohort(w http.ResponseWriter, r *http.Request) {
	// Create a new buffer to read the body, then parse into a []byte
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	uid, err := GetUIDFromBearerToken(r)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorForbidden, w)
		return
	}

	var newCohort model.Cohort

	err = json.Unmarshal(b, &newCohort)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorInternalServer, w)
		return
	}

	c, err := model.AddCohort(newCohort, uid)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorForbidden, w)
		return
	}

	js, err := json.Marshal(c)

	handleErrorAndRespond(js, err, w)
}
