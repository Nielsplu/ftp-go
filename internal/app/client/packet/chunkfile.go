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

type ChunkFileTransfer struct {
	id       string
	fileName string
	fileSize int
}

func performeChunkFile(
	reader *bufio.Reader,
	inChan chan tea.Msg,
	ioErrorChan chan error,
	chunkFileMap *map[string]ChunkFileTransfer,
) {

	// read file path
	strRead, err := reader.ReadString('\n')
	if err != nil {
		ioErrorChan <- err
		return
	}

	// remove last \n from file name
	filePath := strRead[:len(strRead)-1]
	fileName := getFileNameFrom(filePath)

	// read id
	strRead, err = reader.ReadString('\n')
	if err != nil {
		ioErrorChan <- err
		return
	}

	// remove last \n from id
	id := strRead[:len(strRead)-1]

	// read file size
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

	// register filechunk
	(*chunkFileMap)[id] = ChunkFileTransfer{
		id: id, fileSize: fileSize, fileName: fileName,
	}

	// notify iu that a file arrived
	inChan <- view.AddNewProgressBar(fileName)
}

func performeChunk(
	reader *bufio.Reader,
	inChan chan tea.Msg,
	ioErrorChan chan error,
	chunkFileMap *map[string]ChunkFileTransfer,
) {

	// read id
	strRead, err := reader.ReadString('\n')
	if err != nil {
		ioErrorChan <- err
		return
	}

	// remove last \n from id
	id := strRead[:len(strRead)-1]
	chunkFileTransfer, exist := (*chunkFileMap)[id]
	if !exist {
		ioErrorChan <- err
		return
	}

	// read chunk size
	strRead, err = reader.ReadString('\n')
	if err != nil {
		ioErrorChan <- err
		return
	}

	// remove last \n from chunk size
	strFileSize := strRead[:len(strRead)-1]
	chunkSize, err := strconv.Atoi(strFileSize)
	if err != nil {
		inChan <- log.AddLog(log.Sys, "Server responded with a malformed packet")
		return
	}

	destPath := "downloads/" + chunkFileTransfer.fileName
	destFile, err := os.OpenFile(destPath,  os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		ioErrorChan <- err
		return
	}

	defer destFile.Close()

	// read file stat
	statBefore, err := destFile.Stat()
	if err != nil {
		ioErrorChan <- err
		return
	}

	sizeBefore := statBefore.Size()

	// write to file
	bytesCopied, err := io.CopyN(destFile, reader, int64(chunkSize))
	if err != nil || bytesCopied != int64(chunkSize) {
		ioErrorChan <- err
		return
	}

	sizeAfter := sizeBefore + bytesCopied

	var percent float64
	if sizeAfter != int64(chunkFileTransfer.fileSize) {
		percent = float64(sizeAfter) / float64(chunkFileTransfer.fileSize)
	} else {
		inChan <- log.AddLog(log.Sys, "File transfer finished : " + chunkFileTransfer.fileName)
		percent = 1
	}

	inChan <- progressbar.SetProgressFor(chunkFileTransfer.fileName, percent)
}
