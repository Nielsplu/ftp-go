package packet

import (
	"bufio"
	"ftp/internal/app/client/view/cmd"
	"ftp/internal/app/client/view/log"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func getFileNameFrom(path string) string {
	if len(path) < 2 {
		return path
	}

	for i := len(path) - 2; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}

// returns the number of next line
// to print without parsing them
func Performe(
	line string,
	reader *bufio.Reader,
	inChan chan tea.Msg,
	outChan chan string,
	ioErrorChan chan error,
	chunkFileMap *map[string]ChunkFileTransfer,
) int {

	switch line {
	case "Start\n":

		performeStart(reader, inChan, outChan, ioErrorChan)
		return 0

	case "Chunkfile\n":

		performeChunkFile(reader, inChan, ioErrorChan, chunkFileMap)
		return 0

	case "Chunk\n":

		performeChunk(reader, inChan, outChan, ioErrorChan, chunkFileMap)
		return 0

	}

	if strings.HasPrefix(line, "FileCnt") && len(line) > 9 {
		
		fileCntStr := line[8:len(line) - 1]
		lineCnt, err := strconv.Atoi(fileCntStr)
		if err != nil {
			inChan <- log.AddLog(log.Sys, "Server responded with a malformed packet")
			return 0
		}

		inChan <- log.AddLog(log.Srv, line[:len(line) - 1])
		return lineCnt

	}

	if strings.HasPrefix(line, "Moveto") && len(line) > 7{

		inChan <- cmd.ChangeDir(line[7:len(line)-1])
		return 0

	}

	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}

	inChan <- log.AddLog(log.Srv, line)
	return 0
}
