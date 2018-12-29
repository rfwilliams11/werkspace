package handlers

import (
	"testing"

	"bitbucket.centene.com/pdsrtr/rtr-advancement-api-template/env"
)

func TestMain(t *testing.M) {
	env.WorkingDirectory = "../.."
	env.SetConfigOverrideName("base_test")
	t.Run()
}
