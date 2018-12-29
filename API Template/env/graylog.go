package env

import (
	"io"
	"log"
	"os"

 	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"github.com/sirupsen/logrus"
)

func StartLogging(url, facility string){
	if url == "" {
		panic("No Graylog Url for logging.")
	}
	gelfWriter, err := gelf.NewUDPWriter(url)
	if err != nil {
		log.Printf("gelf.NewWriter: %s", err)
	}
	gelfWriter.Facility = facility
	log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
	logrus.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
	log.Print("Logging to Graylog")
}
