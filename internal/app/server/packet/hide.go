package packet

import (
	"fmt"
	"log/slog"
	//"os"

	t "ftp/internal/app/server/types"
)

func PerformeHide(
	path string,
	responseChan chan t.PacketOut,
	state t.ServerState,
) {

	/*if path == "." || path == "/" || len(path) <= 1 || path[len(path)-1] == '\n' {
		responseChan <- t.PacketOut{Buffer: []byte(fmt.Sprintf("File or Directory Unknown : %s\n", path))}
		return
	}*/

	
	_, err := checkPath(state.Config.RootPath, path)
	if err != nil {
        slog.Error(err.Error())
		responseChan <- t.PacketOut{Buffer: []byte(fmt.Sprintf("File or Directory Unknown : %s\n", path))}
		return 
    }

	/*_, err := os.Open(path)
	if err != nil {
        slog.Error(err.Error())
		responseChan <- t.PacketOut{Buffer: []byte(fmt.Sprintf("File or Directory Unknown : %s\n", path))}
		return 
    }*/

	if path[0] == '/' {
		path = path[1:]
	}

	slog.Debug(state.Config.RootPath + path)
	state.HiddenFiles.HidePath(state.Config.RootPath + path)

	responseChan <- t.PacketOut{Buffer: []byte("Ok\n")}

}
