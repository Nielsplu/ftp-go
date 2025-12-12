package types

import "ftp/internal/app/server/hide"

type ServerConfig struct {
	Port, AdminPort, RootPath string
}

type ServerState struct {
	Clients map[string]*Client
	Config  ServerConfig
	HiddenFiles hide.HiddenFileCollection
}