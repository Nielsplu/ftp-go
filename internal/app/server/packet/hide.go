package packet

import (
	"log/slog"

	t "ftp/internal/app/server/types"
)

func PerformeHide(
	path string, 
	responseChan chan t.PacketOut, 
	state t.ServerState,
) {

	if /*path == "." || path == "/" ||*/ len(path) <= 1 || path[len(path) - 1] == '\n' {
		responseChan <- t.PacketOut{Buffer: []byte("File Unknown\n")}
		return
	}

	if path[0] == '/' {
		path = path[1:]
	}

	slog.Debug(state.Config.RootPath + path)
	state.HiddenFiles.HidePath(state.Config.RootPath + path)

	responseChan <- t.PacketOut{Buffer: []byte("Ok\n")}
	
}



