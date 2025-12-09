package remote

import (
	"log/slog"
	"net"
	. "ftp/internal/app/server/types"
)

func ListenOn(
	port *string, 
	admin bool,
	packetInChan chan PacketIn, 
	stopper *Stopper,
) {

	listener, err := net.Listen("tcp", ":" + *port)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("Now listening on port " + *port)

	connChan := make(chan net.Conn)
	errorChan := make(chan error)

	for {

		// start a stopper goroutine to listen, 
		// like that, when shutdwon is required 
		// listener will send error in chanel
		// and StartServer will end
		stopper.Go(func(_ *Stopper) {
			conn, err := listener.Accept()
			if err != nil {
				errorChan <- err
				return
			}
			connChan <- conn
		})

		select {
		case <-stopper.WaitForStopRequest():

			listener.Close()
			stopper.StopChilds()
			slog.Debug("Accept loop shutdown")
			return

		case conn := <- connChan:

			stopper.Go(func (child *Stopper){ 
				// notify app that a new conn as started
				packetInChan <- PacketIn{ Type: NewConn }
				Handle(conn, admin, packetInChan, child) 
			})
		
		case err := <- errorChan:

			slog.Error(err.Error())

		}
	}
}
