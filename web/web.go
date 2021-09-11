package web

import (
	"context"
	"net/http"
	"os"

	"dips/models"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type HostService interface {
	FindHostByFQDN(ctx context.Context, hostname string, domain string) (*models.Host, error)
	CreateHost(ctx context.Context, host *models.Host) error
	DeleateHost(ctx context.Context, host *models.Host) error
}

func InitApp() *http.Server {
	r := mux.NewRouter()
	initRoutes(r)
	return &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8001",
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func initRoutes(r *mux.Router) http.Handler {
	r.Use(loggingMiddleware)
	r.HandleFunc("/api/v1/host/{domain}/{hostname}", GetHostHandler()).Methods("GET")
	r.HandleFunc("/api/v1/host/{domain}/{hostname}", CreateHostHandler()).Methods("POST")
	r.HandleFunc("/api/v1/host/{domain}/{hostname}", DeleteHostHandler()).Methods("DELETE")
	return r
}
