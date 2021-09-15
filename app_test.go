package main

import (
	"os"
	"testing"
)

func TestSendWebHook(t *testing.T) {
	sendWebHook(os.Getenv("FOM_WEBHOOK"), "FOM-OC", "Unit Test", "https://blog.alexellis.io/golang-writing-unit-tests/", "Datefiled", "Messagefiled")
}
