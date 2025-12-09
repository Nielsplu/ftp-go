package packet

import (
	t "ftp/internal/app/server/types"
)

func PerformeList(path string, responseChan chan t.PacketOut) {
	responseChan <- t.PacketOut{ Buffer: []byte(path) }
}