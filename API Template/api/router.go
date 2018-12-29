package api

import (
	"net/http"

	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/api/handlers"
	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/env"
	"github.com/gorilla/mux"
	"github.com/elastic/apm-agent-go/module/apmgorilla"
	log "github.com/sirupsen/logrus"
)

const (
	versionPath = "/version"
	healthPath  = "/health"
)

func router(routeBuilder func(r *mux.Router) *mux.Router) *mux.Router {
	r := mux.NewRouter()

	if env.Config().ApmConfigured() {
		r.Use(apmgorilla.Middleware())
		log.WithFields(log.Fields{"HostURL":env.Config().Apm.ServerUrl,
								  "ServiceName":env.Config().Apm.ServiceName}).Info("Kibana APM agent initialized.")
	} else {
		log.WithFields(log.Fields{"APM Configuration":env.Config().Apm}).Info("Kibana APM agent not configured.")
	}

	r.Use(logger)
	r.Use(swagger)
	r.Use(commonHeaders)

	r.HandleFunc(versionPath, handlers.Version).Methods("GET")
	r.HandleFunc(healthPath, handlers.Health).Methods("GET")
	routeBuilder(r)
	r.PathPrefix("/").Handler(http.StripPrefix("/",fileServer()))
	http.Handle("/", r)

	return r
}

func fileServer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Del("Content-Type")
		http.FileServer(http.Dir(env.WorkingDirectory + "/" + env.Config().PublicDir)).ServeHTTP(w, r)
	})
}

func Default(r *mux.Router) *mux.Router {
	return r
}
