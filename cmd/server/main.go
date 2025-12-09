package main

import (
	"flag"
	"log/slog"
	"ftp/internal/app/server"
)

func parseArgs() (port *string) {

	logLevel := flag.Bool("d", false, "enable debug log level")
	port = flag.String("p", "3333", "server port (default: 3333)")

	flag.Parse()

	if *logLevel {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Set logging level to debug")
	}

	return
}

func main() {
	port := parseArgs()
	server.RunServer(port)
}
