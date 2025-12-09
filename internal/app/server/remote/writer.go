package remote

import (
	. "ftp/internal/app/server/types"
	"log/slog"
	"net"
)

func StartWriter(
	conn net.Conn, 
	packetOutChan chan PacketOut, 
	packetInChan chan PacketIn,
	stopper *Stopper,
) {

	normalQueue := make([][]byte, 0)
	priorityQueue := make([][]byte, 0)

	isChannelFull := false
	nextPacketChan := make(chan []byte, 1)

	for {

		// send to channel condition,
		// may be empty, if not :
		// send to client
		if !isChannelFull {
			if len(priorityQueue) != 0 {
				nextPacketChan <- priorityQueue[0]
				priorityQueue = priorityQueue[1:]
				isChannelFull = true
			} else if len(normalQueue) != 0 {
				nextPacketChan <- normalQueue[0]
				normalQueue = normalQueue[1:]
				isChannelFull = true
			}
		}

		select {
		case <-stopper.WaitForStopRequest():

			slog.Debug("Reader shutdown")
			return

		case packet := <-nextPacketChan:

			isChannelFull = false
			_, err := conn.Write(packet)
			if (err != nil) {
				slog.Error(err.Error())
			}

		case packet := <-packetOutChan:

			if packet.Priority {
				priorityQueue = append(priorityQueue, packet.Buffer)
			} else {
				normalQueue = append(priorityQueue, packet.Buffer)
			}

		}
	}
}
