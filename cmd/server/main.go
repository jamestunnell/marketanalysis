package main

import (
	"os"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jamestunnell/marketanalysis/server"
	"github.com/rs/zerolog"
)

func main() {
	app := kingpin.New("server", "Provide market analysis features with an HTTP server`")
	debug := app.Flag("debug", "Enable debug mode").Bool()
	port := app.Flag("port", "Server port").Required().Int()
	// dir := app.Flag("dir", "Root storage directory (must exist)").Required().String()

	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	srv := server.NewServer(*port)

	srv.Start()

	server.BlockUntilSignaled(syscall.SIGINT, syscall.SIGTERM)

	srv.Stop()
}
