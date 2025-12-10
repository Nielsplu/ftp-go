package packet

import (
	"bufio"
	"ftp/internal/app/client/view/log"

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

func Performe(
	line string,
	reader *bufio.Reader,
	inChan chan tea.Msg,
	ioErrorChan chan error,
	chunkFileMap *map[string]ChunkFileTransfer,
) {

	switch line {
	case "Start\n":

		performeStart(reader, inChan, ioErrorChan)
		return

	case "Chunkfile\n":

		performeChunkFile(reader, inChan, ioErrorChan, chunkFileMap)
		return

	case "Chunk\n":

		performeChunk(reader, inChan, ioErrorChan, chunkFileMap)
		return

	}

	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}

	inChan <- log.AddLog(log.Srv, line)
}
