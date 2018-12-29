package api

import (
	"net/http"
	"time"

	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/env"
	"github.com/gorilla/mux"
)

type Server interface {
	ListenAndServe() error
}

func NewServer(readTimeout time.Duration, writeTimeout time.Duration, routeBuilder func(r *mux.Router) *mux.Router) Server {
	server := &http.Server{
		Handler:      router(routeBuilder),
		Addr:         env.Config().ServerUrl(),
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}
	return server
}
