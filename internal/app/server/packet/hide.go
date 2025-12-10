package packet

import (
	"ftp/internal/app/server/hide"
	"log/slog"

	t "ftp/internal/app/server/types"
)

func PerformeHide(
	path string, 
	responseChan chan t.PacketOut, 
	rootPath string, 
	hiddenFiles hide.HiddenFileCollection,
) {

	if /*path == "." || path == "/" ||*/ len(path) <= 1 || path[len(path) - 1] == '\n' {
		responseChan <- t.PacketOut{Buffer: []byte("File Unknown\n")}
		return
	}

	if path[0] == '/' {
		path = path[1:]
	}

	slog.Debug(rootPath + path)
	hiddenFiles.HidePath(rootPath + path)

	responseChan <- t.PacketOut{Buffer: []byte("Ok\n")}
	
}



