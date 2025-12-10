package server

import "ftp/internal/app/server/hide"

type ServerConfig struct {
	Port, AdminPort, RootPath string
	HiddenFiles hide.HiddenFileCollection
}
