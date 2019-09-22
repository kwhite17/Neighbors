package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/kwhite17/Neighbors/pkg/database"
	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/resources"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

func main() {
	driver := flag.String("dbDriver", "sqlite3", "Name of database driver to use")
	dbHost, dbHostFound := os.LookupEnv("DATABASE_URL")
	developmentMode := flag.Bool("developmentMode", false, "run app in development mode")
	if !dbHostFound {
		dbHost = "file::memory:?mode=memory&cache=shared"
	}
	port, portFound := os.LookupEnv("PORT")
	if !portFound {
		port = "8080"
	}
	flag.Parse()
	log.Println("Connecting to host", dbHost)
	log.Println("Development mode set to:", *developmentMode)

	neighborsDatasource := database.BuildDatasource(*driver, dbHost, *developmentMode)
	shelterManager := &managers.ShelterManager{Datasource: neighborsDatasource}
	itemManager := &managers.ItemManager{Datasource: neighborsDatasource}
	shelterSessionManager := &managers.ShelterSessionManager{Datasource: neighborsDatasource}
	mux := http.NewServeMux()

	mux.Handle("/shelters/", buildShelterServiceHandler(shelterManager, shelterSessionManager, itemManager))
	mux.Handle("/items/", buildItemServiceHandler(shelterSessionManager, itemManager))
	mux.Handle("/session/", buildLoginServiceHandler(shelterManager, shelterSessionManager))
	mux.Handle("/", buildHomeServiceHandler(shelterSessionManager))
	http.ListenAndServe(":"+port, mux)
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

func buildItemServiceHandler(shelterSessionManager *managers.ShelterSessionManager, itemManager *managers.ItemManager) resources.ItemServiceHandler {
	return resources.ItemServiceHandler{
		ItemRetriever:         &retrievers.ItemRetriever{},
		ItemManager:           itemManager,
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
