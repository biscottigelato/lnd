package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/lightningnetwork/lnd/daemon"
	"github.com/roasbeef/btcutil"
)

func main() {
	// Call the "real" main in a nested manner so the defers will properly
	// be executed in the case of a graceful shutdown.
	appDir := btcutil.AppDataDir("lnd", false)
	if err := daemon.LndMain(appDir, nil, nil); err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
