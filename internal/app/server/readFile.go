package server

import (
	"log/slog"
	//"strconv"
	"os"
)

type FileEntry struct {
	filePath string
	fileSize int64
}

func readFile(path string) []FileEntry {

	dir, err := os.Getwd()
	if err != nil {
		slog.Error(err.Error())
	}

	dir += "/../../internal/app/server/"

	data, err := os.ReadDir(dir + path)

	if err != nil {
		slog.Error(err.Error())
		return []FileEntry{}
	}

	var listFichiers []FileEntry

	for _, fichier := range data {

		info, _ := fichier.Info()

		pair := FileEntry{fichier.Name(), info.Size()}

		listFichiers = append(listFichiers, pair)
	}

	return listFichiers
}
