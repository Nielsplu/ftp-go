package client

import (
	"log/slog"
	"net"
)

func Run(remote string) {

	conn , err := net.Dial("tcp", remote)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		conn.Close()
		slog.Debug("Connection closed")
	}()
	slog.Info("Connected to " + conn.RemoteAddr().String())

	
}
