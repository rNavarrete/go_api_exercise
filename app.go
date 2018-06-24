package main

import (
	"database/sql"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//The Initialize method is responsible for creating a database connection and wire up the routes
func (a *App) Initialize(user, password, dbstring string) {

}

//Run method will simply start the application
func (a *App) Run(addr string) {

}
