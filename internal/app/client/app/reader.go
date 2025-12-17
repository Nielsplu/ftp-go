package app

import (
	"bufio"
	"ftp/internal/app/client/packet"
	"ftp/internal/app/client/view/log"
	"ftp/internal/pkg/utils"
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

func StartReader(
	conn 		net.Conn,
	inChan 		chan tea.Msg,
	outChan		chan string,
	ioErrorChan chan error,
	connEndChan chan struct{},
	stopper 	*utils.Stopper,
) {

	reader := bufio.NewReader(conn)
	lineChan := make(chan string, 1)

	chunkFileMap := make(map[string]packet.ChunkFileTransfer, 0)

	nextLineToPrint := 0

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

			if nextLineToPrint > 0 {
				inChan <- log.AddLog(log.Srv, line[:len(line) - 1])
				nextLineToPrint -= 1

				if nextLineToPrint == 0 {
					outChan <- "OK"
				}
				continue
			}

			if line == "End\n" {
				inChan <- tea.Quit()
				connEndChan <- struct{}{}
				println("Connection terminée par le serveur")
				return
			}
			
			nextLineToPrint = packet.Performe(line, reader, inChan, outChan, ioErrorChan, &chunkFileMap)
		}
	}
}
