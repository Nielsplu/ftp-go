package packet

import t "ftp/internal/app/server/types"

func PerformeCd(
	path string, 
	responseChan chan t.PacketOut, 
	state *t.ServerState,
	client *t.Client,
) {

	fileChecked, err := checkPath(state, client.CurrentPath, path)
	if err != nil || !fileChecked.isDir || fileChecked.isHidden {
		responseChan <- t.PacketOut{ Buffer: []byte("FileUnknown\n") }
		return
	}

	client.CurrentPath = fileChecked.clientVisiblePath
	responseChan <- t.PacketOut{ 
		Buffer: []byte("Moveto " + fileChecked.clientVisiblePath + "\n"),
	}

}