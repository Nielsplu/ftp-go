package packet

import (
	t "ftp/internal/app/server/types"
)

func PerformeHide(
	path string, 
	responseChan chan t.PacketOut, 
	state *t.ServerState,
	client *t.Client,
) {

	fileChecked, err := checkPath(state, client.CurrentPath, path)
	if err != nil || fileChecked.isHidden {
		responseChan <- t.PacketOut{Buffer: []byte("FileUnknown\n")}
		return 
	}

	state.HiddenFiles.HidePath(fileChecked.clientVisiblePath)
	responseChan <- t.PacketOut{ Buffer: []byte("Ok\n") }
	
}
