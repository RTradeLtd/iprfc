package iprfc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	pbr "github.com/RTradeLtd/grpc/lens/request"

	ipfsapi "github.com/RTradeLtd/go-ipfs-api"
	"github.com/RTradeLtd/iprfc/lens"
)

var (
	// error is returned when we've downloaded the last rfc
	errMoreRFCs = errors.New("no more rfcs to download")
	baseURL     = "https://tools.ietf.org/pdf/"
	// https://tools.ietf.org/pdf/rfc5245.pdf
)

// GetRFC gets an RFC number
func GetRFC(num int) string {
	return fmt.Sprintf("rfc%v", num)
}

// FormatURL returns a url to download an RFC
func FormatURL(rfc string) string {
	return baseURL + rfc + ".pdf"
}

// GetAndSave is used to download an RFC as a PFD
func GetAndSave(rfc string) error {
	url := FormatURL(rfc)
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
	return ioutil.WriteFile(rfc+".pdf", body, os.FileMode(0640))
}

// DownloadAndSave is used to download and save a file
func DownloadAndSave(max int) {
	var count = 1
	for {
	START:
		// max of 0 mens no more to download
		// this allows us to do testing without downloading everything
		if max != 0 && count > max {
			return
		}
		err := GetAndSave(GetRFC(count))
		switch err {
		case nil:
			count++
			goto START
		case errMoreRFCs:
			count++
			goto START
		default:
			log.Fatalf("error downloading rfc: %s", err)
		}
	}
}

// StoreAndIndex is used to store a file on IPFS and index it
//
// It reads all files in the current directory, adds it to IPFS, and then indexing it against Lens
func StoreAndIndex(ctx context.Context, sh *ipfsapi.Shell, lc *lens.Client, index bool) error {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".pdf") {
			fh, err := os.Open(file.Name())
			if err != nil {
				return err
			}
			hash, err := sh.Add(fh)
			if err != nil {
				return err
			}
			fmt.Printf("added\t%s\t%s\n", hash, file.Name())
			if index {
				if err := Index(ctx, lc, hash); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Index is used to index a hash against lens
func Index(ctx context.Context, lc *lens.Client, hash string) error {
	_, err := lc.Index(ctx, &pbr.Index{
		Type:       "ipld",
		Identifier: hash,
	})
	return err
}
