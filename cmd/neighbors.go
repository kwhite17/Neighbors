package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kwhite17/Neighbors/pkg/database"
	"github.com/kwhite17/Neighbors/pkg/items"
	"github.com/kwhite17/Neighbors/pkg/login"
	"github.com/kwhite17/Neighbors/pkg/shelters"
)

var NeighborsDatabase database.NeighborsDatasource

func main() {
	dbHost := flag.String("dbhost", "file::memory:?mode=memory&cache=shared", "Name of host on which to run Neighbors")
	developmentMode := flag.Bool("developmentMode", false, "run app in development mode")
	flag.Parse()
	log.Println("Connecting to host", *dbHost)
	log.Println("Development mode set to:", *developmentMode)

	NeighborsDatabase = database.NeighborsDatasource{Database: database.InitDatabase(*dbHost, *developmentMode)}
	shelterManager := &shelters.ShelterManager{Datasource: NeighborsDatabase}
	shelterSessionManager := &login.ShelterSessionManager{Datasource: NeighborsDatabase}
	mux := http.NewServeMux()

	mux.Handle("/shelters/", buildShelterServiceHandler(shelterManager, shelterSessionManager))
	mux.Handle("/items/", buildItemServiceHandler())
	mux.Handle("/login/", buildLoginServiceHandler(shelterManager, shelterSessionManager))
	http.ListenAndServe(":8080", mux)
}

func buildShelterServiceHandler(shelterManager *shelters.ShelterManager, shelterSessionManager *login.ShelterSessionManager) shelters.ShelterServiceHandler {
	return shelters.ShelterServiceHandler{
		ShelterRetriever:      &shelters.ShelterRetriever{},
		ShelterManager:        shelterManager,
		ShelterSessionManager: shelterSessionManager,
	}
}

func buildItemServiceHandler() items.ItemServiceHandler {
	return items.ItemServiceHandler{
		ItemRetriever: &items.ItemRetriever{},
		ItemManager:   &items.ItemManager{Datasource: NeighborsDatabase},
	}
}

func buildLoginServiceHandler(shelterManager *shelters.ShelterManager, shelterSessionManager *login.ShelterSessionManager) login.LoginServiceHandler {
	return login.LoginServiceHandler{
		ShelterSessionManager: shelterSessionManager,
		LoginRetriever:        &login.LoginRetriever{},
	}
}
