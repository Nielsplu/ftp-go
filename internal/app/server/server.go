package server

import (
	"log/slog"
	"net"
)

func StartServer(port *string, stopper *Stopper) {

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("Now listening on port " + *port)

	connChan := make(chan net.Conn)
	errorChan := make(chan error)

	for {

		// start a stopper goroutine to listen, 
		// like that when shutdwon is required 
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
		case <-stopper.Wait():

			slog.Debug("Stopping accept loop ...")
			listener.Close()
			stopper.Stop()
			return

		case conn := <- connChan:

			stopper.Go(func (child *Stopper){ 
				Handle(conn, child) 
			})
		
		case err := <- errorChan:
			slog.Error(err.Error())
		}
	}
}
