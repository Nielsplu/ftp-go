package packet

import t "ftp/internal/app/server/types"

func PerformeCd(
	path string, 
	responseChan chan t.PacketOut, 
	state t.ServerState,
	client *t.Client,
) {

	if len(path) == 0 {
		client.CurrentPath = "/"
		responseChan <- t.PacketOut{ Buffer: []byte("Moveto /\n") }
		return
	}

	if path[0] != '/' {
		path = client.CurrentPath + "/" + path
	}

	cleanPath, err := checkPath(state.Config.RootPath, path)
	if err != nil {
		responseChan <- t.PacketOut{ Buffer: []byte("FileUnknown\n") }
		return
	}

	client.CurrentPath = cleanPath
	responseChan <- t.PacketOut{ Buffer: []byte("Moveto " + cleanPath + "\n") }

}