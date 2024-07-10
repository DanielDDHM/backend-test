package internal

import (
	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/internal/handlers"
)

type App struct {
	logger *charmLog.Logger
}

func NewApp(logger *charmLog.Logger) *App {
	return &App{
		logger: logger,
	}
}

func (a *App) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/breeds", handlers.GetBreeds).Methods("GET")
	r.HandleFunc("/breeds/{id}", handlers.GetBreed).Methods("GET")
	r.HandleFunc("/breeds", handlers.CreateBreed).Methods("POST")
	r.HandleFunc("/breeds/{id}", handlers.UpdateBreed).Methods("PUT")
	r.HandleFunc("/breeds/{id}", handlers.DeleteBreed).Methods("DELETE")
}