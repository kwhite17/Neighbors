package main

import (
	"net/http"

	"github.com/kwhite17/Neighbors/pkg/database"

	"github.com/kwhite17/Neighbors/pkg/shelters"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle("/shelters/", buildShelterServiceHandler())

	http.ListenAndServe(":8080", mux)
}

func buildShelterServiceHandler() shelters.ShelterServiceHandler {
	return shelters.ShelterServiceHandler{
		ShelterRetriever: &shelters.ShelterRetriever{},
		ShelterManager:   &shelters.ShelterManager{Datasource: database.NeighborsDatabase},
	}
}
