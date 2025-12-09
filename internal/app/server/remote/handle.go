package remote

import (
	"bufio"
	. "ftp/internal/app/server/types"
	"ftp/internal/app/server/packet"
	"log/slog"
	"net"
)

func Handle(
	conn net.Conn, 
	admin bool,
	packetInChan chan PacketIn, 
	stopper *Stopper,
) {

	defer func() {
		conn.Close()
		slog.Debug("Conn close")
	}()

	connCloseChan := make(chan struct{})
	
	packetOutChan := make(chan PacketOut, 10)
	stopper.Go(func(child *Stopper) {
		// when writer ends, reader must too 
		defer func (){ connCloseChan <- struct{}{} }()
		StartWriter(conn, packetOutChan, packetInChan, child)
	})

	reader := bufio.NewReader(conn)
	lineChan := make(chan string, 1)

	for {

		stopper.Go(func(_ *Stopper) {
			line, err := reader.ReadString('\n')
			if (err != nil) {
				// if read error, conn is close
				slog.Error(err.Error())
				connCloseChan <- struct{}{}
				return
			}
			lineChan <- line
		})

		select {
			case <-stopper.WaitForStopRequest():

				stopper.StopChilds()
				slog.Debug("Reader shutdown")
				return

			case <- connCloseChan:

				stopper.StopChilds()
				slog.Debug("Reader shutdown")

				// notify app that conn is closed
				packetInChan <- PacketIn { Type: ConnEnd }
				return
				
			case line := <- lineChan:

				packetIn, err := packet.Parse(line)
				if (err != nil) {
					slog.Error(err.Error())
					continue
				}

				if (packetIn.Type == End) {
					connCloseChan <- struct{}{}
					return
				}

				packetInChan <- packetIn
		}


	}

}