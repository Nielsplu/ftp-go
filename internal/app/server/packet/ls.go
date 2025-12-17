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
	isDir    bool
}

func PerformeList(
	path string, 
	responseChan chan t.PacketOut, 
	state *t.ServerState,
	client *t.Client,
) {

	checkFile, err := checkPath(state, client.CurrentPath, path)
	if err != nil || !checkFile.isDir || checkFile.isHidden {
		responseChan <- t.PacketOut{ Buffer: []byte("FileUnknown\n") }
		return
	}

	arrayFileEntry := readDir(checkFile.absPath, checkFile.clientVisiblePath, state.HiddenFiles)
	fileCnt := len(arrayFileEntry)

	packetStr := fmt.Sprintf("FileCnt %d\n", fileCnt)

	// format packet
	for _, entry := range arrayFileEntry {
		if entry.isDir {
			packetStr += fmt.Sprintf("%s/\n", entry.filePath)
		}else {
			packetStr += fmt.Sprintf("%s %d\n", entry.filePath, entry.fileSize)
		}
	}

	responseChan <- t.PacketOut{ Buffer: []byte(packetStr) }
}

func readDir(
	dirAbsPath, clientVisibleDirPath string, 
	hiddenFiles hide.HiddenFileCollection,
) []fileEntry {

	// open the file and get data
	data, err := os.ReadDir(dirAbsPath)
	if err != nil {
		slog.Error(err.Error())
		return []fileEntry{}
	}

	// fill arrayDirectory
	var fileEntries []fileEntry
	for _, file := range data {

		info, _ := file.Info()
		if !hiddenFiles.IsPathHidden(clientVisibleDirPath + "/" + file.Name()){
			fileEntries = append(fileEntries, fileEntry{ file.Name(), info.Size(), file.IsDir() })
		}

	}

	return fileEntries
}
