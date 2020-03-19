package main

import (
	"flag"

	"github.com/RTradeLtd/iprfc"
)

var (
	max         = flag.Int("max.rfc", 1, "the maximum rfc to download")

)

func init() {
	flag.Parse()
}

func main() {
	iprfc.DownloadAndSave(*max)
}