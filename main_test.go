package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestIPRFC(t *testing.T) {
	t.Cleanup(func() {
		os.RemoveAll("downloads")
	})
	if err := os.Mkdir("downloads", os.FileMode(0640)); err != nil {
		t.Fatal(err)
	}
	downloadAndSave(1)
	data, err := ioutil.ReadFile("downloads/rfc1.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if data == nil {
		t.Fatal("no data found")
	}
}
