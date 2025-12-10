package app

import (
	"bufio"
	"ftp/internal/app/client/packet"
	"ftp/internal/pkg/utils"
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

func StartReader(
	conn 		net.Conn,
	inChan 		chan tea.Msg,
	ioErrorChan chan error,
	connEndChan chan struct{},
	stopper 	*utils.Stopper,
) {

	reader := bufio.NewReader(conn)
	lineChan := make(chan string, 1)

	for {

		stopper.Go(func(_ *utils.Stopper) {
			line, err := reader.ReadString('\n')
			if err != nil {
				ioErrorChan <- err
				return
			}
				
			lineChan <- line
		})

		select {
		case <-stopper.WaitForStopRequest():

			stopper.StopChilds()
			return
		
		case line := <-lineChan:

			if line == "End\n" {
				inChan <- tea.Quit()
				connEndChan <- struct{}{}
				println("Connection terminée par le serveur")
				return
			}

			packet.Performe(line, reader, inChan, ioErrorChan)
		}
	}
}
