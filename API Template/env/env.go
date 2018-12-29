package env

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	instance         Conf
	once             sync.Once
	baseName         = "base"
	splitCred 		 DecryptionClient
	WorkingDirectory = func() string {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(fmt.Errorf("Unable to get working directory: %s\n", err))
		}
		return dir
	}()
)

func Initialize(credentials DecryptionClient){
	splitCred = credentials
	instance = newConfig()
}

func Config() Conf {
	return instance
}

func newConfig() Conf {
	c := Conf{}
	initEnvTags(reflect.TypeOf(c), "")

	viper.SetConfigType("yaml")
	viper.AddConfigPath(WorkingDirectory + "/config/")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName(baseName)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Unable to read config file %s: %s\n", baseName, err))
	}

	if err := viper.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("Unable to unmarshal config: %s\n", err))
	}

	c.Mongo.Password = loadFromSplitPea(c.Splitpea.Decrypt, c.Mongo.EUUID, c.Mongo.MongoKey1, c.Mongo.MongoKey2, c.Splitpea.Username, c.Splitpea.Password)
	c.Mongo.Url = buildUrl(c.Mongo.Username, c.Mongo.Password, c.Mongo.Url)

	return c
}

// This is what builds the splitpea request and either gets a successful response or panics if there is an error.
func loadFromSplitPea(url, euuid, key1, key2, username, password string) string {

	req := SplitPeaRequest{euuid, []string{key1, key2}}
	res, err := splitCred.Call(url, username, password, req)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("SplitPea Decrypt Error")
	}

	return res
}

// This is what binds the 'env' struct tag to the associated struct parameter
// Generally, adding any new types should be as simple as adding it to the first case
func initEnvTags(t reflect.Type, parent string) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		switch f.Type.Kind() {
		case reflect.String, reflect.Int64:
			if env := strings.ToUpper(f.Tag.Get("env")); env != "" {
				// The next two lines are both required for viper to recognize and bind the env var
				viper.BindEnv(strings.ToLower(fmt.Sprintf("%s%s", parent, f.Name)), env)
				viper.SetDefault(env, getDefault(f.Type))
			}
		case reflect.Struct:
			initEnvTags(f.Type, fmt.Sprintf("%s%s.", parent, f.Name))
		default:
			if f.Tag.Get("env") != "" {
				log.Infof("No case for type [%v] to bind config value to env variable", f.Type)
			}
		}
	}
}

func getDefault(t reflect.Type) interface{} {
	zero := reflect.Zero(t)
	switch t.Kind() {
	case reflect.String:
		return zero.String()
	case reflect.Int64:
		return zero.Int()
	default:
		log.Info("No default for type [%v] for config", t)
		return nil
	}
}

func buildUrl(username string, password string, url string) string {
	newUrl := "mongodb://" + username + ":" + password + "@" + url + "/PDM?authSource=admin"
	return newUrl
}

