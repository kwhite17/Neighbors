package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kwhite17/Neighbors/pkg/database"
	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/resources"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

var NeighborsDatabase database.NeighborsDatasource

func main() {
	driver := flag.String("dbDriver", "sqlite3", "Name of database driver to use")
	dbHost := flag.String("dbhost", "file::memory:?mode=memory&cache=shared", "Name of host on which to run Neighbors")
	developmentMode := flag.Bool("developmentMode", false, "run app in development mode")
	flag.Parse()
	log.Println("Connecting to host", *dbHost)
	log.Println("Development mode set to:", *developmentMode)

	NeighborsDatabase = database.NeighborsDatasource{Database: database.InitDatabase(database.BuildConfig(*driver, *dbHost, *developmentMode))}
	shelterManager := &managers.ShelterManager{Datasource: NeighborsDatabase}
	itemManager := &managers.ItemManager{Datasource: NeighborsDatabase}
	shelterSessionManager := &managers.ShelterSessionManager{Datasource: NeighborsDatabase}
	mux := http.NewServeMux()

	mux.Handle("/shelters/", buildShelterServiceHandler(shelterManager, shelterSessionManager, itemManager))
	mux.Handle("/items/", buildItemServiceHandler(shelterSessionManager))
	mux.Handle("/session/", buildLoginServiceHandler(shelterManager, shelterSessionManager))
	mux.Handle("/", buildHomeServiceHandler(shelterSessionManager))
	//mux.HandleFunc("/", loadHomePage)
	http.ListenAndServe(":8080", mux)
}

func buildHomeServiceHandler(shelterSessionManager *managers.ShelterSessionManager) resources.HomeServiceHandler {
	return resources.HomeServiceHandler{
		ShelterSessionManager: shelterSessionManager,
	}
}

func buildShelterServiceHandler(shelterManager *managers.ShelterManager, shelterSessionManager *managers.ShelterSessionManager, itemManager *managers.ItemManager) resources.ShelterServiceHandler {
	return resources.ShelterServiceHandler{
		ShelterRetriever:      &retrievers.ShelterRetriever{},
		ShelterManager:        shelterManager,
		ShelterSessionManager: shelterSessionManager,
		ItemManager:           itemManager,
	}
}

func buildItemServiceHandler(shelterSessionManager *managers.ShelterSessionManager) resources.ItemServiceHandler {
	return resources.ItemServiceHandler{
		ItemRetriever:         &retrievers.ItemRetriever{},
		ItemManager:           &managers.ItemManager{Datasource: NeighborsDatabase},
		ShelterSessionManager: shelterSessionManager,
	}
}

func buildLoginServiceHandler(shelterManager *managers.ShelterManager, shelterSessionManager *managers.ShelterSessionManager) resources.LoginServiceHandler {
	return resources.LoginServiceHandler{
		ShelterManager:        shelterManager,
		ShelterSessionManager: shelterSessionManager,
		LoginRetriever:        &retrievers.LoginRetriever{},
	}
}
