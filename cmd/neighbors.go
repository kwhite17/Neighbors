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
	userManager := &managers.UserManager{Datasource: neighborsDatasource}
	itemManager := &managers.ItemManager{Datasource: neighborsDatasource}
	userSessionManager := &managers.UserSessionManager{Datasource: neighborsDatasource}

	mux := http.NewServeMux()
	mux.Handle("/shelters/", buildUserServiceHandler(userManager, userSessionManager, itemManager))
	mux.Handle("/items/", buildItemServiceHandler(userSessionManager, itemManager))
	mux.Handle("/session/", buildLoginServiceHandler(userManager, userSessionManager))
	mux.Handle("/", buildHomeServiceHandler(userSessionManager))
	http.ListenAndServe(":"+port, mux)
}

func buildHomeServiceHandler(userSessionManager *managers.UserSessionManager) resources.HomeServiceHandler {
	return resources.HomeServiceHandler{
		UserSessionManager: userSessionManager,
	}
}

func buildUserServiceHandler(userManager *managers.UserManager, userSessionManager *managers.UserSessionManager, itemManager *managers.ItemManager) resources.UserServiceHandler {
	return resources.UserServiceHandler{
		UserRetriever:      &retrievers.ShelterRetriever{},
		UserManager:        userManager,
		UserSessionManager: userSessionManager,
		ItemManager:        itemManager,
	}
}

func buildItemServiceHandler(userSessionManager *managers.UserSessionManager, itemManager *managers.ItemManager) resources.ItemServiceHandler {
	return resources.ItemServiceHandler{
		ItemRetriever:      &retrievers.ItemRetriever{},
		ItemManager:        itemManager,
		UserSessionManager: userSessionManager,
	}
}

func buildLoginServiceHandler(userManager *managers.UserManager, userSessionManager *managers.UserSessionManager) resources.LoginServiceHandler {
	return resources.LoginServiceHandler{
		UserManager:        userManager,
		UserSessionManager: userSessionManager,
		LoginRetriever:     &retrievers.LoginRetriever{},
	}
}
