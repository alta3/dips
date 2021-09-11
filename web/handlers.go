package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"dips/models"

	"github.com/gorilla/mux"
)

func urlToFQDN(r *http.Request) (models.FQDN, error) {
	vars := mux.Vars(r)
	varDomain := vars["domain"]
	varHostname := vars["hostname"]
	return models.FQDN{Domain: varDomain, Hostname: varHostname}, nil
}

func CreateHostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fqdn, err := urlToFQDN(r)
		if err != nil {
			log.Printf("Cannot parse domain/host from url. err=%v \n", err)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		host, err := models.CreateHost(fqdn)
		if err != nil {
			log.Printf("Cannot create new Host. err=%v \n", err)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		hostJSON, err := json.Marshal(host)
		if err != nil {
			log.Printf("Cannot marshal Host into json. err=%v \n", err)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Printf("%+v\n", host)
		w.Write(hostJSON)
	}
}

func GetHostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fqdn, err := urlToFQDN(r)
		if err != nil {
			log.Printf("Cannot parse domain/host from url. err=%v \n", err)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		host, err := models.FindHostByFQDN(fqdn)
		if err != nil {
			log.Printf("Host not in DB. err=%v \n", err)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		hostJSON, err := json.Marshal(host)
		if err != nil {
			log.Printf("Cannot marshal Host into json. err=%v \n", err)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(hostJSON)
	}
}

func DeleteHostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fqdn, err := urlToFQDN(r)
		if err != nil {
			log.Printf("Cannot parse domain/host from url. err=%v \n", err)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = models.DeleteHostByFQDN(fqdn)
		if err != nil {
			log.Printf("Host delete failed. err=%v \n", err)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
