package remote

import (
	"bufio"
	"ftp/internal/app/server/packet"
	t "ftp/internal/app/server/types"
	"ftp/internal/pkg/utils"
	"log/slog"
	"net"
)

func Handle(
	conn net.Conn,
	admin bool,
	packetInChan chan t.PacketIn,
	stopper *utils.Stopper,
) bool {

	// Do not defer conn.Close here
	// it's writer responsability to 
	// close the connection, like 
	// it needs to send END\n before

	// channel used to coordinate writer 
	// and reader when an error occurs
	// set capacity to one to avoid blocking
	ioErrorChan := make(chan error, 1)

	// create writer goroutine with priority
	packetOutChan := make(chan t.PacketOut, 10)
	stopper.Go(func(child *utils.Stopper) {
		
		err := StartWriter(conn, packetOutChan, packetInChan, child)

		// send error if reader 
		// hasn't send one yet
		if len(ioErrorChan) < cap(ioErrorChan) {
			ioErrorChan <- err
		}

	})

	reader := bufio.NewReader(conn)
	lineChan := make(chan string, 1)

	for {

		// create a goroutine to read like
		// we need to listen for stopper.
		// If stop is requested, conn is closed and
		// err sent to ioErrorChan will be ignored
		stopper.Go(func(_ *utils.Stopper) {
			line, err := reader.ReadString('\n')
			if err != nil {

				// send error if writer 
				// hasn't send one yet
				if len(ioErrorChan) != 1 { 
					ioErrorChan <- err
				}

				return
			}

			lineChan <- line
		})

		select {
		case <-stopper.WaitForStopRequest():

			// stop writer and reader
			stopper.StopChilds()
			slog.Debug("Conn shutdown (order)")
			return true

		case err := <-ioErrorChan:

			// stop writer and reader
			stopper.StopChilds()
			slog.Error("Conn shutdown with error : " + err.Error())

			// notify core that conn has ended
			packetInChan <- t.PacketIn { Type: t.ConnEnd }
			return false

		case line := <-lineChan:

			packetIn, err := packet.Parse(line, packetOutChan, admin)
			if err != nil {
				slog.Error("Parse error : " + err.Error())
				continue
			}

			if packetIn.Type == t.End {

				// stop writer and reader
				stopper.StopChilds()
				slog.Debug("Conn shutdown (End requested)")

				// notify core that conn has ended
				packetInChan <- t.PacketIn { Type: t.ConnEnd }
				return false
			}

			packetInChan <- packetIn
		}
	}
}
