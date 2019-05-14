package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kwhite17/Neighbors/pkg/database"
	"github.com/kwhite17/Neighbors/pkg/items"
	"github.com/kwhite17/Neighbors/pkg/shelters"
)

var NeighborsDatabase database.NeighborsDatasource

func main() {
	dbHost := flag.String("dbhost", ":memory:", "Name of host on which to run Neighbors")
	developmentMode := flag.Bool("developmentMode", false, "run app in development mode")
	flag.Parse()
	log.Println("Connecting to host", *dbHost)
	log.Println("Development mode set to:", *developmentMode)

	NeighborsDatabase = database.NeighborsDatasource{Database: database.InitDatabase(*dbHost, *developmentMode)}
	mux := http.NewServeMux()

	mux.Handle("/shelters/", buildShelterServiceHandler())
	mux.Handle("/items/", buildItemServiceHandler())
	http.ListenAndServe(":8080", mux)
}

func buildShelterServiceHandler() shelters.ShelterServiceHandler {
	return shelters.ShelterServiceHandler{
		ShelterRetriever: &shelters.ShelterRetriever{},
		ShelterManager:   &shelters.ShelterManager{Datasource: NeighborsDatabase},
	}
}

func buildItemServiceHandler() items.ItemServiceHandler {
	return items.ItemServiceHandler{
		ItemRetriever: &items.ItemRetriever{},
		ItemManager:   &items.ItemManager{Datasource: NeighborsDatabase},
	}
}
