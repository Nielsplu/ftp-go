package packet

import t "ftp/internal/app/server/types"

func PerformeCd(
	path string, 
	responseChan chan t.PacketOut, 
	rootPath string,
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

	cleanPath, err := checkPath(rootPath, path)
	if err != nil {
		responseChan <- t.PacketOut{ Buffer: []byte("FileUnknown\n") }
		return
	}

	clientVisiblePath := cleanPath[len(rootPath):]

	client.CurrentPath = clientVisiblePath
	responseChan <- t.PacketOut{ Buffer: []byte("Moveto " + clientVisiblePath + "\n") }

}