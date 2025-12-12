package packet

import (
	"fmt"
	t "ftp/internal/app/server/types"
	"ftp/internal/pkg/utils"
	"io"
	"log/slog"
	"os"
)

func PerformeGet(
	path string,
	responseChan chan t.PacketOut,
	state *t.ServerState,
	client *t.Client,
	stopper *utils.Stopper,
) {

	fileChecked, err := checkPath(state, client.CurrentPath, path)
	if err != nil || fileChecked.isDir || fileChecked.isHidden {
		responseChan <- t.PacketOut{ Buffer: []byte("FileUnknown\n") }
		return
	}

	//open file
	file, err := os.Open(fileChecked.absPath)
	if err != nil {
		responseChan <- t.PacketOut{Buffer: []byte("FileUnknown\n")}
		return
	}

	//file info
	info, err := file.Stat()
	if err != nil {
		responseChan <- t.PacketOut{Buffer: []byte("FileUnknown\n")}
		file.Close()
		return
	}

	//choose the transfer style
	if info.Size() < 20*1_000_000 {
		sendInOnePacket(file, info, path, responseChan)
	} else {
		sendWithChunks(file, info, path, responseChan, stopper)
	}

}

func sendInOnePacket(file *os.File, info os.FileInfo, path string, responseChan chan t.PacketOut) {
	defer file.Close()

	slog.Debug("sending with start")
	header := fmt.Sprintf(
		"Start\n"+
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
	responseChan <- t.PacketOut{Buffer: packet}
}

func sendWithChunks(
	file *os.File,
	info os.FileInfo,
	path string,
	responseChan chan t.PacketOut,
	stopper *utils.Stopper,
) {
	id := utils.MakeId()

	// initialize file transfer
	initPacket := fmt.Sprintf(
		"Chunkfile\n"+
			"%s\n"+
			"%s\n"+
			"%d\n",
		path, id, info.Size(),
	)

	responseChan <- t.PacketOut{Buffer: []byte(initPacket)}

	slog.Debug("starting to send file chunk by chunk")

	//  chunks

	stopper.Go(func(child *utils.Stopper) {
		defer file.Close()

		for {

			// check if stop is required
			select {
			case <-child.WaitForStopRequest():
				return
			default:
			}

			buffer := make([]byte, 4096)

			bytesRead, err := file.Read(buffer)
			if err != nil {
				if err != io.EOF {
					slog.Error("Error during read file : "+path, "erreur", err)
				}
				break
			}

			header := fmt.Sprintf(
				"Chunk\n"+
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
		}
	})
}
