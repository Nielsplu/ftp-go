package server

import (
	"log/slog"
	//"strconv"
	"os"
	"bufio"
	"fmt"
)

var dir string

type FileEntry struct {
	filePath string
	fileSize int64
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

func CreateArrayOfHidenFile(/*file string*/) /*[]string*/{
	dir, err := os.Getwd()
	if err != nil {
		slog.Error(err.Error())
	}

	dir += "/../../internal/app/server/"

	file, err := os.Open(dir + "hidden_file.txt")
	if err != nil {
        slog.Error(err.Error())
		return
    }

    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(line)
    }

	fmt.Println("here")
    if err := scanner.Err(); err != nil {
        slog.Error(err.Error())
		return
    }
	
}
