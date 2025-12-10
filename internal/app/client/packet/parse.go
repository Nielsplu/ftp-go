package packet

import (
	"bufio"
	"ftp/internal/app/client/view"
	"ftp/internal/app/client/view/log"
	"ftp/internal/app/client/view/progressbar"
	"io"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func getFileNameFrom(path string) string {
	if len(path) < 2 { return path }

	for i := len(path) - 2; i >= 0; i -- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}

func performeStart(
	reader *bufio.Reader,
	inChan chan tea.Msg,
	ioErrorChan chan error,
) {
	strRead, err := reader.ReadString('\n')
	if err != nil {
		ioErrorChan <- err
		return
	}

	// remove last \n from file name
	filePath := strRead[:len(strRead) - 1]
	fileName := getFileNameFrom(filePath)

	// notify iu that a file arrived
	inChan <- view.AddNewProgressBar(filePath)

	strRead, err = reader.ReadString('\n')
	if err != nil {
		ioErrorChan <- err
		return
	}

	// remove last \n from file size
	strFileSize := strRead[:len(strRead) - 1]
	fileSize, err := strconv.Atoi(strFileSize)
	if err != nil {
		inChan <- log.AddLog(log.Sys, "Server responded with a malformed packet")
		return
	}

	destFile, err := os.Create("data/" + fileName)
	if err != nil {
		inChan <- log.AddLog(log.Sys, "Erreur lors de l'ouverture du fichier local : " + fileName)
		return
	}
	defer destFile.Close()

	bytesCopied, err := io.CopyN(destFile, reader, int64(fileSize))
	if err != nil || bytesCopied != int64(fileSize) {
		inChan <- log.AddLog(log.Sys, "Erreur lors du teléchargement du fichier : " + filePath)
		return
	}

	inChan <- progressbar.SetProgressFor(filePath, 1)
}

func Performe(
	line   string, 
	reader *bufio.Reader,
	inChan chan tea.Msg,
	ioErrorChan chan error,
) {
	if line == "Start\n" {
		performeStart(reader, inChan, ioErrorChan)
		return
	}

	if len(line) > 0 && line[len(line) - 1] == '\n' {
		line = line[:len(line) - 1]
	}

	inChan <- log.AddLog(log.Srv, line)
}