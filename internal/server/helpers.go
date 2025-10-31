package server

import (
	"awesomeProject1/internal/server/builders"
	"awesomeProject1/internal/store"
	"awesomeProject1/pkg/logger"
	"github.com/gorilla/mux"
)

func (a *Api) configureRouter() {
	router := mux.NewRouter()
	builder := &builders.Builder{
		Router: router,
		Store:  a.store,
	}

	builders.BuildRoutes(builder)
	a.router = router
}

func (a *Api) configureLogger() error {
	return logger.Init(a.config.LogLevel, a.config.LogDir)
}

func (a *Api) configureStore() error {
	store := store.New(a.config.Storage)

	if err := store.Start(); err != nil {
		return err
	}
	a.store = store
	return nil
}
