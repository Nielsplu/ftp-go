package remote

import (
	"fmt"
	t "ftp/internal/app/server/types"
	"ftp/internal/pkg/utils"
	"log/slog"
	"net"
)

func ListenOn(
	port *string, 
	admin bool,
	packetInChan chan t.PacketIn, 
	stopper *utils.Stopper,
) {

	listener, err := net.Listen("tcp", ":" + *port)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("Now listening on port " + *port)

	connChan := make(chan net.Conn)
	errorChan := make(chan error, 1)

	for {

		// start a stopper goroutine to listen, 
		// like that, when shutdwon is required 
		// and listener closed, an error 
		// will be send in chanel
		// and this goroutine will end
		stopper.Go(func(_ *utils.Stopper) {
			conn, err := listener.Accept()
			if err != nil {
				errorChan <- err
				return
			}
			connChan <- conn
		})

		select {
		case <-stopper.WaitForStopRequest():

			// close accept goroutine
			listener.Close()

			// wait for all conn to shutdown
			stopper.StopChilds()

			slog.Debug(fmt.Sprintf("Accept loop shutdown (order, admin=%v)", admin))
			return

		case conn := <- connChan:

			// do not create a goroutine
			// if in admin mode
			if (!admin) {

				stopper.Go(func(child *utils.Stopper) { 
					// notify app that a new conn as started
					packetInChan <- t.PacketIn{ Type: t.NewConn }
					Handle(conn, admin, packetInChan, child) 
				})

			}else {

				// handle connection without goroutine to 
				// avoid multiple admin connection
				stopRequested := Handle(conn, admin, packetInChan, stopper) 
				if stopRequested { 
					// close accept goroutine
					listener.Close()

					slog.Debug(fmt.Sprintf("Accept loop shutdown (order, admin=%v)", admin))
					return
				}

			}
		
		case err := <- errorChan:

			slog.Error("Accept loop error : " + err.Error())

		}
	}
}