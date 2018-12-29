package rtrapi

import (
	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/api"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/env"
	gelf "github.com/seatgeek/logrus-gelf-formatter"

	"time"
)

func Start(routeBuilder func(*mux.Router) *mux.Router) {
	StartWithTimeouts(15 * time.Second, 15 * time.Second, routeBuilder)
}
func StartWithTimeouts(readTimeout time.Duration, writeTimeout time.Duration, routeBuilder func(*mux.Router) *mux.Router) {
	srv := api.NewServer(readTimeout, writeTimeout, routeBuilder)
	level := env.Config().GetLogLevel()
	log.SetLevel(level)
	graylogUrl := env.Config().GetGraylogUrl()
	log.SetFormatter(new(gelf.GelfFormatter))
	env.StartLogging(graylogUrl, env.Config().LogFacility)

	log.Info("Starting Server...")
	log.WithFields(log.Fields{"Version": env.Version}).Info()
	log.WithFields(log.Fields{"Branch": env.Branch}).Info()
	log.WithFields(log.Fields{"SHA1": env.SHA1}).Info()
	log.WithFields(log.Fields{"Build Time": env.BuildTime}).Info()

	log.WithFields(log.Fields{"MongoDB Connection": env.Config().MongoUrl()}).Info()

	log.WithField("Log level", level).Debug()
	log.WithField("Log facility", env.Config().LogFacility).Info()
	log.WithFields(log.Fields{"Listening on": env.Config().ServerUrl()}).Info()
	log.Fatal(srv.ListenAndServe())
}
