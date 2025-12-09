package main

import (
	"flag"
	"log/slog"
	"ftp/internal/app/server"
)

func parseArgs() (port *string, adminPort *string) {

	logLevel := flag.Bool("debug", false, "enable debug log level")
	port = flag.String("port", "3333", "server port (default: 3333)")
	adminPort = flag.String("admin-port", "4444", "server admin port (default: 4444)")

	flag.Parse()

	if *logLevel {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Set logging level to debug")
	}

	return
}

func main() {
	port, adminPort := parseArgs()
	server.Start(server.ServerConfig {
		Port: *port, AdminPort: *adminPort,
	})
}
