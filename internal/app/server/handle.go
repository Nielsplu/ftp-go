package server

import (
	"log/slog"
	"net"
)

func Handle(conn net.Conn, stopper *Stopper) {
	defer func() {
		conn.Close()
		slog.Debug("Connection close")
	}()


}