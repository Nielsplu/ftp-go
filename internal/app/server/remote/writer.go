package remote

import (
	t "ftp/internal/app/server/types"
	"ftp/internal/pkg/utils"
	"net"
)

func StartWriter(
	conn net.Conn, 
	packetOutChan chan t.PacketOut, 
	packetInChan chan t.PacketIn,
	stopper *utils.Stopper,
	resetTimerChan chan struct{},
) error {

	defer conn.Close()

	lowPriorityQueue := make([][]byte, 0)
	normalQueue := make([][]byte, 0)

	isChannelFull := false
	nextPacketChan := make(chan []byte, 1)

	for {

		// send to channel condition,
		// may be empty, if not :
		// send to client
		if !isChannelFull {
			if len(normalQueue) != 0 {
				nextPacketChan <- normalQueue[0]
				normalQueue = normalQueue[1:]
				isChannelFull = true
			} else if len(lowPriorityQueue) != 0 {
				nextPacketChan <- lowPriorityQueue[0]
				lowPriorityQueue = lowPriorityQueue[1:]
				isChannelFull = true
			}
		}

		select {
		case <-stopper.WaitForStopRequest():

			// send end before closing the conn
			conn.Write([]byte("END\n"))
			return nil

		case packet := <-nextPacketChan:

			isChannelFull = false
			_, err := conn.Write(packet)
			if (err != nil) {
				return err
			}

			resetTimerChan <- struct{}{}

		case packet := <-packetOutChan:

			if packet.LowPriority {
				lowPriorityQueue = append(lowPriorityQueue, packet.Buffer)
			} else {
				normalQueue = append(normalQueue, packet.Buffer)
			}

		}
	}
}
