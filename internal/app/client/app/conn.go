package app

import (
	"ftp/internal/app/client/view"
	"ftp/internal/pkg/utils"
	"net"
	
	tea "github.com/charmbracelet/bubbletea"
)

func Handle(
	outChan chan string,
	inChan chan tea.Msg,
	conn net.Conn,
) {

	ioErrorChan := make(chan error, 10)
	stopper := utils.NewStopper()

	stopper.Go(func(child *utils.Stopper) {
		StartReader(conn, inChan, ioErrorChan, stopper.WaitForStopRequest(), child)
	})

	stopper.Go(func(child *utils.Stopper) {
		StartWriter(conn, outChan, inChan, ioErrorChan, stopper.WaitForStopRequest(), child)
	})

	for {
		select {
		case <-stopper.WaitForStopRequest():

			// garcefull stop
			conn.Close()
			stopper.StopChilds()
			inChan <- tea.Quit()
			return

		case err := <-ioErrorChan:

			// error stop
			inChan <- view.NotifyFatalError("Connection perdu", err)
			conn.Close()
			stopper.StopChilds()
			return

		}
	}
}
