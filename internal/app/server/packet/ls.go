package packet

import (
	t "ftp/internal/app/server/types"
	"os"
	"log/slog"
)

var dir string

type FileEntry struct {
	filePath string
	fileSize int64
}

func PerformeList(path string, responseChan chan t.PacketOut) {
	responseChan <- t.PacketOut{ Buffer: []byte(path) }
}

func readFile(path string) []FileEntry {

	data, err := os.ReadDir(dir + path)

	if err != nil {

		slog.Error(err.Error())
		return []FileEntry{}
	}

	var arrayDirectory []FileEntry

	for _, fichier := range data {

		info, _ := fichier.Info()

		fileEntry := FileEntry{fichier.Name(), info.Size()}

		arrayDirectory = append(arrayDirectory, fileEntry)
	}

	return arrayDirectory
}

