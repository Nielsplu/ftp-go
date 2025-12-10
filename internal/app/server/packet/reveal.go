package packet

import (
	t "ftp/internal/app/server/types"
	"ftp/internal/app/server/hide"
)

func PerformeReveal(
	path string, 
	responseChan chan t.PacketOut, 
	rootPath string,
	hiddenFiles hide.HiddenFileCollection,
) {
	if !hiddenFiles.IsPathHidden(path){
		responseChan <- t.PacketOut{Buffer: []byte("File Not Hidden\n")}
	}

	if path[0] == '/' {
		path = path[1:]
	}
	
	hiddenFiles.RevealPath(rootPath + path)

	responseChan <- t.PacketOut{Buffer: []byte("Ok\n")}
}