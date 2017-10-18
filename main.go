//go:generate go-bindata -pkg data -o pkg/data/data.go data/
package main

import (
	"net/http"

	"github.com/cairesvs/beeru/pkg/database"
	"github.com/cairesvs/beeru/pkg/logger"

	"github.com/cairesvs/beeru/pkg/router"
	"github.com/gorilla/mux"
)

func main() {
	db := database.GetInstance()
	db.LoadToDatabase()
	r := mux.NewRouter()
	r.HandleFunc("/pdv/{id:[0-9]+}", router.GetPDV).Methods("GET")
	r.HandleFunc("/pdv", router.CreatePDV).Methods("POST")
	r.HandleFunc("/pdvs", router.FindPDV).Methods("GET")
	logger.Info("Running on :8000")
	logger.Fatal(http.ListenAndServe(":8000", r))
}
