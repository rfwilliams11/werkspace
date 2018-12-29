package env

import (
	"os"
	"testing"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"reflect"
	"github.com/spf13/viper"
	"fmt"
)

var (
	mongoKey1      = "MONGO_KEY1"
	mongoKey2      = "MONGO_KEY2"
	splitPeaPassword = "SPLITPEA_PASSWORD"

	environentKeys  = []string{mongoKey1, mongoKey2, splitPeaPassword}
)

type splitCredMock struct {
	responseKey string
}

func (s splitCredMock) Call(url string, username string, password string, query SplitPeaRequest) (response string, err error) {
	response = s.responseKey
	err = nil
	return
}

// What is the point of this?
func initialize() Conf {

	// Reset Viper TODO: Is this actually necessary?
	viper.Reset()

	// Initialize the environment with the the provided configuration
	splitPeaResponse := "mongopassword"
	Initialize(splitCredMock{splitPeaResponse})

	baseConfigPath := fmt.Sprintf("../config/%v.yml", baseName)
	baseBytes, _ := ioutil.ReadFile(baseConfigPath)
	expected := Conf{}
	yaml.Unmarshal([]byte(strings.ToLower(string(baseBytes))), &expected)

	keyOne, keyOneSet := os.LookupEnv(mongoKey1)
	if keyOneSet {
		expected.Mongo.MongoKey1 = keyOne
	}

	keyTwo, keyTwoSet := os.LookupEnv(mongoKey2)
	if keyTwoSet {
		expected.Mongo.MongoKey2 = keyTwo
		expected.Mongo.Password = splitPeaResponse
	}

	spPass, spPassSet := os.LookupEnv(splitPeaPassword)
	if spPassSet {
		expected.Splitpea.Password = spPass
	}

	// We need to correct the mongo URL which is built using the mongo data from configuration
	expected.Mongo.Url = "mongodb://" + expected.Mongo.Username + ":" + expected.Mongo.Password + "@" + expected.Mongo.Url + "/PDM?authSource=admin"

	return expected
}

func resetEnvironmentVariables(){
	for _, key := range environentKeys {
		os.Unsetenv(key)
	}
}

func TestMain(t *testing.M) {
	baseName = "base_test"
	WorkingDirectory = ".."
	resetEnvironmentVariables()
	t.Run()
}

func TestEnvBaseConfig(t *testing.T) {
	expected := initialize()

	if actual := Config(); !reflect.DeepEqual(expected, actual) {
		t.Errorf("Base configuration not read correctly, expected %+v, got %+v", expected, actual)
	}
}

func TestEnvBaseConfigEnvironmentSecrets(t *testing.T) {
	os.Setenv(mongoKey1, "mongokeyone")
	os.Setenv(mongoKey2, "mongokeytwo")
	os.Setenv(splitPeaPassword, "passwordistaco")
	expected := initialize()

	if actual := Config(); !reflect.DeepEqual(expected, actual) {
		t.Errorf("Base configuration not read correctly, expected %+v, got %+v", expected, actual)
	}
}

func TestEnvEmptyConfig(t *testing.T) {
	baseName = "empty_test"
	resetEnvironmentVariables()
	expected := initialize()

	if actual := Config(); !reflect.DeepEqual(expected, actual) {
		t.Errorf("Config not read correctly, expected %+v, got %+v", expected, actual)
	}
}