package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kwhite17/Neighbors/pkg/database"
	"github.com/kwhite17/Neighbors/pkg/items"
	"github.com/kwhite17/Neighbors/pkg/samaritans"

	"github.com/kwhite17/Neighbors/pkg/neighbors"
)

func main() {
	directory, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	directory, err = filepath.EvalSymlinks(directory)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(directory)
	mux := http.NewServeMux()
	mux.Handle("/neighbors/", neighbors.NeighborServiceHandler{Database: database.NeighborsDatabase})
	mux.Handle("/samaritans/", samaritans.SamaritanServiceHandler{Database: database.NeighborsDatabase})
	mux.Handle("/items/", items.ItemServiceHandler{Database: database.NeighborsDatabase})
	mux.Handle("/templates/", http.FileServer(http.Dir(directory+"/templates/")))

	http.ListenAndServe(":8080", mux)
}
