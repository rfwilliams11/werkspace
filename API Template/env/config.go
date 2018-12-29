package env

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// This struct should match the format of the config yaml files
// Env variable overrides are set with the 'env' tag
// '.' are replaced by '_' when viper is searching for the env variable
// By default, the env variable bound is the uppercase parameter name
// Viper delimits nested parameters with a '.'
type Conf struct {
	Server struct {
		Host string
		Port string
	}
	Apm struct {
		ServerUrl	string `env:"ELASTIC_APM_SERVER_URL"`
		ServiceName	string `env:"ELASTIC_APM_SERVICE_NAME"`
		Environment	string `env:"ELASTIC_APM_ENVIRONMENT"`
		Version	string 	   `env:"ELASTIC_APM_SERVICE_VERSION"`
	}
	Mongo struct {
		Url         string `env:"MONGO_URL"`
		DialTimeout int64  `env:"MONGO_DIAL_TIMEOUT"`
		DbName      string `env:"MONGO_DB_NAME"`
		Username    string `env:"MONGO_USERNAME"`
		Password    string
		MongoKey1   string `env:"MONGO_KEY1"`
        MongoKey2   string `env:"MONGO_KEY2"`
        EUUID       string `env:"SPLITPEA_EUUID"`
	}
	Splitpea struct {
		Username string `env:"SPLITPEA_USERNAME"`
		Password string `env:"SPLITPEA_PASSWORD"`
		Decrypt  string `env:"SPLITPEA_DECRYPT"`
	}
	SwaggerEndpoint    string `env:"SWAGGER_ENDPOINT"`
	PublicDir          string `env:"PUBLIC_DIR"`
	LogFacility    	   string `env:"LOG_FACILITY"`
	LogLevel           string `env:"LOG_LEVEL"`
	MongoSecretPath    string `env:"MONGO_SECRET_PATH"`
	SplitpeaSecretPath string `env:"SPLITPEA_SECRET_PATH"`
	GraylogUrl		   string `env:"GRAYLOG_URL"`
}

func (c Conf) ServerUrl() string {
	return c.Server.Host + ":" + c.Server.Port
}

func (c Conf) MongoUrl() string {
	url := c.Mongo.Url
	urlComponents := strings.SplitAfter(url, "@")
	if len(urlComponents) > 1 {
		url = urlComponents[len(urlComponents)-1]
	} else {
		url = "mongoDB"
	}
	return url
}

func (c Conf) GetLogLevel() log.Level {
	l, e := log.ParseLevel(strings.ToLower(c.LogLevel))
	if e != nil {
		panic(e)
	}
	return l
}

func (c Conf) GetGraylogUrl() string {
	return c.GraylogUrl
}

func (c Conf) ApmConfigured() bool {
	return len(c.Apm.ServerUrl) > 0
}

