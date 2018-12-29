package api

import (
	"bytes"
	"fmt"
	"net/http"

	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/env"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.NewV4()

		log.WithFields(log.Fields{
			"RequestId": id,
			"Method": r.Method,
			"Request": r.URL,
			"Headers": printHeader(r.Header)}).Debug()

		wrapper := &responseWrapper{w, http.StatusOK}
		next.ServeHTTP(wrapper, r)

		responseFields := log.Fields{
			"RequestId": id,
			"Method": r.Method,
			"Request": r.URL,
			"ResponseCode": wrapper.status,
			"Headers": printHeader(wrapper.w.Header())}

		if wrapper.status >= 400 {
			log.WithFields(responseFields).Warn()
		} else {
			log.WithFields(responseFields).Debug()
		}
	})
}

func swagger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, exists := r.URL.Query()["url"]; !exists && r.URL.String() == "/api/" {
			scheme := r.URL.Scheme
			if scheme == "" {
				scheme = "http"
			}
			swagger := fmt.Sprintf("%v://%v%v", scheme, r.Host, env.Config().SwaggerEndpoint)
			http.Redirect(w, r, fmt.Sprintf("%v://%s/api/?url=%s", scheme, r.Host, swagger), http.StatusMovedPermanently)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Expires", "0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "https://confluence.centene.com")
		next.ServeHTTP(w, r)
	})
}

func printHeader(h http.Header) string {
	buffer := bytes.NewBufferString("[")
	for k, v := range h {
		buffer.WriteString(fmt.Sprintf("%s:%v,", k, v))
	}
	if buffer.Len() > 1 {
		buffer.Truncate(buffer.Len() - 1)
	}
	buffer.WriteString("]")
	return buffer.String()
}

type responseWrapper struct {
	w      http.ResponseWriter
	status int
}

func (r *responseWrapper) Header() http.Header {
	return r.w.Header()
}

func (r *responseWrapper) Write(b []byte) (int, error) {
	return r.w.Write(b)
}

func (r *responseWrapper) WriteHeader(c int) {
	r.status = c
	r.w.WriteHeader(c)
}
