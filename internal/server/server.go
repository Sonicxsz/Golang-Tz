package server

import (
	"awesomeProject1/internal/store"
	"awesomeProject1/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
	config *Config
	router *mux.Router
	store  *store.Store
}

func New(config *Config) *Api {
	return &Api{
		config: config,
	}
}

func (api *Api) Start() error {
	if err := api.configureLogger(); err != nil {
		return err
	}
	defer logger.Log.Close()

	if err := api.configureStore(); err != nil {
		return err
	}

	api.configureRouter()

	return http.ListenAndServe(api.config.BindAddr, api.router)
}
