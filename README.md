# IPRFC

`iprfc` is a tool to download all RFCs in PDF form, store them on IPFS, and index them using the Lens search engine.

# Installation

Before proceeding you'll need to have a valid install of Go 1.14 to build. 

1) Download dependencies with `go mod download`
2) Build with `make` and an executable called `iprfc` will be created in the current directory

# Usage

```
NAME:
   iprfc - a tool to download all known RFCs in PDF and add them to IPFS

USAGE:
   iprfc [global options] command [command options] [arguments...]

DESCRIPTION:
   It requires at a minimum being able to access a go-ipfs node, and optionally a Lens endpoint to index against

COMMANDS:
   download-and-save  download all known RFCs and save
   store-and-index    store RFCs onto IPFS and index
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --max.rfc value        the maximum rfc to download, 0 means no max (default: 1)
   --ipfs.endpoint value  the go-ipfs api endpoint to use (default: "127.0.0.1:5001")
   --lens.endpoint value  the lens grpc endpoint to use (default: "127.0.0.1:9998")
   --index                whether or not to initiate lens indexing (default: false)
   --help, -h             show help (default: false)
```

## Downloading RFCs

The most basic functionality of this tool consists of download all available RFCs in PDF format, saving them in the current directory. To prevent you from accidentally spamming the IETF website, the default setting is to download the first RFC, and then exit. This can be configured with the `--max.rfc` flag.

To download the first 2 RFCs run `iprfc --max.rfc 2 download-and-save`. 

To download all available RFCs run `iprfc --max.rfc 0 download-save`. Note that this will require you manually exit the process. I couldn't think of a good way to detect when finished downloading all RFCs. Initially I tried using the 404 status code, but apparently some RFC numbers dont exist, and this turned out to not be a good way. PRs welcomed for this functionality.

## Storing On IPFS And Indexing

Before doing this you'll need to download the RFCs either using the `download-and-save` command, or doing it manually, but who wants to do it manually? ;)

Because not everyone will have access to a Lens gRPC endpoint (Lens is open-source btw so you can easily do this), the default behavior of the `store-and-index` command is simply to add the RFCs to IPFS. 

One thing to note is that this will pick up **ANY** PDF's in the current directory, so make sure you run this without any sensitive files in place.

To save all RFCs onto ipfs and not index run `iprfc store-and-index`.

To save all RFCs onto IPFS and index run `iprfc --index store-and-index`.