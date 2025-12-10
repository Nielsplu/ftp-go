package packet

import (
	"fmt"
	t "ftp/internal/app/server/types"
	"log/slog"
	"os"
)

type fileEntry struct {
	filePath string
	fileSize int64
}

func PerformeList(path string, responseChan chan t.PacketOut, rootPath string) {
	arrayFileEntry := readFile(rootPath + "/" + path)
	fileCnt := len(arrayFileEntry)

	packetStr := fmt.Sprintf("FileCnt %d\n", fileCnt)

	for _, entry := range arrayFileEntry {
		packetStr += fmt.Sprintf("%s %d\n", entry.filePath, entry.fileSize)
	}

	responseChan <- t.PacketOut{Buffer: []byte(packetStr)}
}

func readFile(path string) []fileEntry {

	data, err := os.ReadDir(path)

	if err != nil {

		slog.Error(err.Error())
		return []fileEntry{}
	}

	var arrayDirectory []fileEntry

	for _, fichier := range data {

		info, _ := fichier.Info()

		fileEntry := fileEntry{fichier.Name(), info.Size()}

		arrayDirectory = append(arrayDirectory, fileEntry)
	}

	return arrayDirectory
}
