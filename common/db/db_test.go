package db

import (
	"os"
	"testing"
)

func Test_EmptyConnString(t *testing.T) {
	os.Unsetenv("CCHC_DBSTR")
	_, err := getConnString()
	if err == nil {
		t.Error("Did not get an error when DB connection string was not set")
	}
}

func Test_SetConnString(t *testing.T) {
	connstr := "postgress://user:pass@localhost:5432/cchc"
	envvar := "CCHC_DBSTR"
	os.Setenv(envvar, connstr)
	got, err := getConnString()
	if err != nil {
		t.Error("Got an error when connection string was set: ", err)
	}
	if got != connstr {
		t.Error("Connection string does not match environment variable")
	}
}
