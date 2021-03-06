package iprfc

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestIPRFC(t *testing.T) {
	t.Cleanup(func() {
		os.Remove("rfc1.pdf")
	})
	DownloadAndSave(1)
	data, err := ioutil.ReadFile("rfc1.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if data == nil {
		t.Fatal("no data found")
	}
}
