package packet

import (
	t "ftp/internal/app/server/types"
)

func PerformeReveal(
	path string, 
	responseChan chan t.PacketOut, 
	state t.ServerState,
) {
	if !state.HiddenFiles.IsPathHidden(path){
		responseChan <- t.PacketOut{Buffer: []byte("File Not Hidden\n")}
	}

	if path[0] == '/' {
		path = path[1:]
	}
	
	state.HiddenFiles.RevealPath(state.Config.RootPath + path)

	responseChan <- t.PacketOut{Buffer: []byte("Ok\n")}
}