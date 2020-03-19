package main

import (
	"context"
	"log"
	"os"

	"github.com/RTradeLtd/config"
	ipfsapi "github.com/RTradeLtd/go-ipfs-api"
	"github.com/RTradeLtd/iprfc"
	"github.com/RTradeLtd/iprfc/lens"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "iprfc"
	app.Usage = "a tool to download all known RFCs in PDF and add them to IPFS"
	app.Description = "It requires at a minimum being able to access a go-ipfs node, and optionally a Lens endpoint to index against"
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:  "max.rfc",
			Usage: "the maximum rfc to download, 0 means no max",
			Value: 1,
		},
		&cli.StringFlag{
			Name:  "ipfs.endpoint",
			Usage: "the go-ipfs api endpoint to use",
			Value: "127.0.0.1:5001",
		},
		&cli.StringFlag{
			Name:  "lens.endpoint",
			Usage: "the lens grpc endpoint to use",
			Value: "127.0.0.1:9998",
		},
		&cli.BoolFlag{
			Name:  "index",
			Usage: "whether or not to initiate lens indexing",
			Value: false,
		},
	}
	app.Commands = cli.Commands{
		{
			Name:        "download-and-save",
			Usage:       "download all known RFCs and save",
			Description: "this will download all known RFCs and save to the current directory",
			Action: func(c *cli.Context) error {
				iprfc.DownloadAndSave(c.Int("max.rfc"))
				return nil
			},
		},
		{
			Name:        "store-and-index",
			Usage:       "store RFCs onto IPFS and index, default behavior is to just add to IPFS",
			Description: "this uses go-ipfs-api and the Lens gRPC client to add RFC pdfs to IPFS and index them",
			Action: func(c *cli.Context) error {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				sh := ipfsapi.NewShell(c.String("ipfs.endpoint"))
				var lcc *lens.Client
				if c.Bool("index") {
					cfg := config.Endpoints{}
					cfg.Lens.URL = c.String("lens.endpoint")
					lc, err := lens.NewClient(cfg)
					if err != nil {
						return err
					}
					lcc = lc
				}
				return iprfc.StoreAndIndex(ctx, sh, lcc, c.Bool("index"))
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
