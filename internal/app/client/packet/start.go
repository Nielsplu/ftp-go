package packet

import (
	"bufio"
	"fmt"
	"ftp/internal/app/client/view"
	"ftp/internal/app/client/view/log"
	"ftp/internal/app/client/view/progressbar"
	"io"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

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
	filePath := strRead[:len(strRead)-1]
	fileName := getFileNameFrom(filePath)

	// notify iu that a file arrived
	inChan <- view.AddNewProgressBar(filePath)

	strRead, err = reader.ReadString('\n')
	if err != nil {
		ioErrorChan <- err
		return
	}

	// remove last \n from file size
	strFileSize := strRead[:len(strRead)-1]
	fileSize, err := strconv.Atoi(strFileSize)
	if err != nil {
		inChan <- log.AddLog(log.Sys, "Server responded with a malformed packet")
		return
	}

	destFile, err := os.OpenFile("downloads/" + fileName, os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		inChan <- log.AddLog(log.Sys, "Erreur lors de l'ouverture du fichier local : " + fileName)
		return
	}
	defer destFile.Close()

	bytesCopied, err := io.CopyN(destFile, reader, int64(fileSize))
	if err != nil || bytesCopied != int64(fileSize) {
		inChan <- log.AddLog(log.Sys, "Erreur lors du teléchargement du fichier : "+filePath)
		inChan <- log.AddLog(log.Sys, fmt.Sprintf("%d et %d", ))
		return
	}

	inChan <- progressbar.SetProgressFor(filePath, 1)
}