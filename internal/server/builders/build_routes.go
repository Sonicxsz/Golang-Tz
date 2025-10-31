package builders

import (
	"awesomeProject1/internal/handlers"
	"awesomeProject1/internal/service"
	"awesomeProject1/internal/store"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	url = "/api/v1"
)

type Builder struct {
	Router *mux.Router
	Store  *store.Store
}

func BuildRoutes(b *Builder) {
	//Subscriptions
	subscriptionService := service.NewSubscriptionService(b.Store.SubscriptionRepository())
	subscriptionHandler := handlers.NewCatalogHandler(subscriptionService)
	b.Router.HandleFunc(url+"/subscription", subscriptionHandler.Create).Methods("POST")
	b.Router.HandleFunc(url+"/subscription", subscriptionHandler.Update).Methods("PATCH")
	b.Router.HandleFunc(url+"/subscription/{id}", subscriptionHandler.Delete).Methods("DELETE")
	b.Router.HandleFunc(url+"/subscription/{id}", subscriptionHandler.GetById).Methods("GET")
	b.Router.HandleFunc(url+"/subscriptions/total", subscriptionHandler.GetTotal).Methods("GET")
	b.Router.HandleFunc(url+"/subscriptions", subscriptionHandler.GetAll).Methods("GET")
	// Swagger UI
	b.Router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
