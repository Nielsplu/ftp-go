package packet

import (
	"fmt"
	"ftp/internal/app/server/hide"
	t "ftp/internal/app/server/types"
	"log/slog"
	"os"
)

type fileEntry struct {
	filePath string
	fileSize int64
}

func PerformeList(
	path string, 
	responseChan chan t.PacketOut, 
	state t.ServerState,
) {

	arrayFileEntry := readFile(state.Config.RootPath + path, state.HiddenFiles)
	fileCnt := len(arrayFileEntry)

	packetStr := fmt.Sprintf("FileCnt %d\n", fileCnt)

	// add line per line the Name 
	// and the Size of the file
	for _, entry := range arrayFileEntry {
		packetStr += fmt.Sprintf("%s %d\n", entry.filePath, entry.fileSize)
	}

	responseChan <- t.PacketOut{Buffer: []byte(packetStr)}
}

func readFile(path string, hiddenFiles hide.HiddenFileCollection) []fileEntry {

	// open the file and get data
	data, err := os.ReadDir(path)

	if err != nil {

		slog.Error(err.Error())
		return []fileEntry{}
	}

	var arrayDirectory []fileEntry

	if path[len(path) - 1] != '/'{
		path += "/"
	}

	// add to arrayDirectory the Name and the Size 
	// of the file not in hiddenFiles
	for _, fichier := range data {

		info, _ := fichier.Info()

		if !hiddenFiles.IsPathHidden(path + fichier.Name()){

			fileEntry := fileEntry{fichier.Name(), info.Size()}

			arrayDirectory = append(arrayDirectory, fileEntry)
		}
	}

	return arrayDirectory
}
