// Package router is the http controller for PDV functionalities.
package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/cairesvs/beeru/pkg/database"
	"github.com/cairesvs/beeru/pkg/model"
	"github.com/gorilla/mux"
)

type genericErr struct {
	Message string `json:"message,omitempty"`
}

type successResponse struct {
	Message string `json:"message,omitempty"`
}

func jsonError(msg string) []byte {
	b, _ := json.Marshal(&genericErr{msg})
	return b
}

func jsonResponse(w http.ResponseWriter, v interface{}, status int) {
	b, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonError(fmt.Sprintf("Error marshalling %s", err)))
	} else {
		w.WriteHeader(status)
		w.Write(b)
	}
}

// GetPDV returns pdv by id
func GetPDV(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	be := database.GetInstance()
	input := &database.PDVGetInput{
		Database: be,
		ID:       vars["id"],
	}
	pdv := database.GetPDV(input)
	if pdv == nil {
		jsonResponse(w, pdv, http.StatusNotFound)
	} else {
		jsonResponse(w, pdv, http.StatusNotFound)
	}
}

// CreatePDV creates a pdv and return the pdv created
func CreatePDV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p *model.PDV
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonError(fmt.Sprintf("%s", err)))
		return
	}
	input := &database.PDVCreateInput{
		Database: database.GetInstance(),
		PDV:      p,
	}
	_, err = database.CreatePDV(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonError(fmt.Sprintf("%s", err)))
		return
	}
	jsonResponse(w, &successResponse{"The PDV was inserted"}, http.StatusOK)
}

func FindPDV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	be := database.GetInstance()
	lat, err := strconv.ParseFloat(r.FormValue("lat"), 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonError(fmt.Sprintf("%s", err)))
		return
	}
	lng, err := strconv.ParseFloat(r.FormValue("lng"), 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonError(fmt.Sprintf("%s", err)))
		return
	}
	input := &database.PDVFindInput{
		Database: be,
		Point: &model.Point{
			Type:        "Point",
			Coordinates: []float64{lat, lng},
		},
	}
	pdvs := database.FindPDV(input)
	jsonResponse(w, pdvs, http.StatusOK)
}
