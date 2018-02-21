package main

import (
	"net/http"

	"github.com/kwhite17/Neighbors/database"

	"github.com/kwhite17/Neighbors/neighbors"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/neighbors/", neighbors.NeighborServiceHandler{Database: database.NeighborsDatabase})
	// mux.HandleFunc("/samaritans/", samaritans.RequestHandler)
	// mux.Handle("/items", items.ItemsServiceHandler{})
	http.ListenAndServe(":8080", mux)
}
