package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	// error is returned when we've downloaded the last rfc
	errMoreRFCs = errors.New("no more rfcs to download")
	baseURL     = "https://tools.ietf.org/pdf/"
	max         = flag.Int("max.rfc", 1, "the maximum rfc to download")
	// https://tools.ietf.org/pdf/rfc5245.pdf
)

func init() {
	flag.Parse()
}

func getRFC(num int) string {
	return fmt.Sprintf("rfc%v", num)
}

// formatURL returns a url to download an RFC
func formatURL(rfc string) string {
	return baseURL + rfc + ".pdf"
}

// GetAndSave is used to download an RFC as a PFD
func GetAndSave(rfc string) error {
	url := formatURL(rfc)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return errMoreRFCs
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("downloads/"+rfc+".pdf", body, os.FileMode(0640))
}

func downloadAndSave(max int) {
	var count = 1
	for {
	START:
		// max of 0 mens no more to download
		// this allows us to do testing without downloading everything
		if max != 0 && count > max {
			return
		}
		err := GetAndSave(getRFC(count))
		switch err {
		case nil:
			count++
			goto START
		case errMoreRFCs:
			log.Println("finished downloading rfc")
		default:
			log.Fatalf("error downloading rfc: %s", err)
		}
	}
}

func main() {
	if err := os.MkdirAll("downloads", os.FileMode(0640)); err != nil {
		log.Fatal(err)
	}
	downloadAndSave(0)
}
