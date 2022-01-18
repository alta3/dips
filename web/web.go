package web

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Config struct {
	Network string  // 10.0.0.0/12
	Gateway string  // 10.0.0.1
	DhcpStartAddress string // 10.0.2.1
	DhcpEndAddress  string // 10.15.255.254
	DhcpLease string // 8h
	ListenIP string // 0.0.0.0
	ListenPort string //8001
}

var conf Config

func InitApp(c Config) *mux.Router {
	conf = c
	r := mux.NewRouter()
	initRoutes(r)
	return r
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
