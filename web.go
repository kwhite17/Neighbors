package main

import (
	"net/http"

	"github.com/kwhite17/Neighbors/items"
	"github.com/kwhite17/Neighbors/neighbors"
	"github.com/kwhite17/Neighbors/samaritans"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/neighbors/", neighbors.RequestHandler)
	mux.HandleFunc("/samaritans/", samaritans.RequestHandler)
	mux.Handle("/items", items.ItemsServiceHandler{})
	http.ListenAndServe(":8080", mux)
}
