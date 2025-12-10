package packet

import (
	"fmt"
	t "ftp/internal/app/server/types"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func makeId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func PerformeGet(
	path string,
	responseChan chan t.PacketOut,
	rootPath string,
) {

	go func() {
		//get the absolute path
		absRoot, err := filepath.Abs(rootPath)
		if err != nil {
			slog.Error("Invalid root path")
			responseChan <- t.PacketOut{Buffer: []byte("FileUnknown\n")}
			return
		}

		//get the path
		fullPath := filepath.Join(absRoot, path)

		//clean path
		cleanPath := filepath.Clean(fullPath)

		if !strings.HasPrefix(cleanPath, absRoot) {
			slog.Warn("error path traversal attempt", "path", path)
			responseChan <- t.PacketOut{Buffer: []byte("FileUnknown\n")}
			return
		}

		//open file
		file, err := os.Open(cleanPath)
		if err != nil {
			responseChan <- t.PacketOut{Buffer: []byte("FileUnknown\n")}
			return
		}
		defer file.Close()

		//file info
		info, err := file.Stat()
		if err != nil || info.IsDir() {
			responseChan <- t.PacketOut{Buffer: []byte("FileUnknown\n")}
			return
		}

		//choose the transfer style
		if info.Size() < 20*1_000_000 {
			sendInOnePacket(file, info, path, responseChan)
		} else {
			sendWithChunks(file, info, path, responseChan)
		}
	}()

}

func sendInOnePacket(file *os.File, info os.FileInfo, path string, responseChan chan t.PacketOut) {

	header := fmt.Sprintf(
		"START\n"+
			"%s\n"+
			"%d\n",
		path, info.Size(),
	)

	headerBytes := []byte(header)

	// get the all size file
	buffer := make([]byte, info.Size())

	// read all file
	_, err := io.ReadFull(file, buffer)
	if err != nil {
		slog.Error("Error reading entire file", "err", err)
		return
	}

	packet := append(headerBytes, buffer...)
	responseChan <- t.PacketOut{ Buffer: packet }
}

func sendWithChunks(file *os.File, info os.FileInfo, path string, responseChan chan t.PacketOut) {
	id := makeId()

	// initialize file transfer
	initPacket := fmt.Sprintf(
		"START\n"+
			"%s\n"+
			"%s\n"+
			"%d\n",
		path, id, info.Size(),
	)

	responseChan <- t.PacketOut{Buffer: []byte(initPacket)}

	//  chunks
	for {

		buffer := make([]byte, 4096)

		bytesRead, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			slog.Error("Error during read file : "+path, "erreur", err)
			break
		}

		if bytesRead == 0 {
			if err == io.EOF {
				break
			}
			continue
		}

		header := fmt.Sprintf(
			"CHUNK\n"+
				"%s\n"+
				"%d\n",
			id, bytesRead,
		)

		headerBytes := []byte(header)
		packet := append(headerBytes, buffer[:bytesRead]...)

		responseChan <- t.PacketOut{
			Buffer:      packet,
			LowPriority: true,
		}

		if err == io.EOF {
			break
		}

	}
}
