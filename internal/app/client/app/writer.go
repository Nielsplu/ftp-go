package app

import (
	"ftp/internal/app/client/view"
	"ftp/internal/app/client/view/log"
	"ftp/internal/app/client/view/progressbar"
	"ftp/internal/pkg/utils"
	"net"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func StartWriter(
	conn net.Conn,
	outChan chan string,
	inChan chan tea.Msg,
	ioErrorChan chan error,
	connEndChan chan struct{},
	stopper *utils.Stopper,
) {

	for {
		select {
		case <-stopper.WaitForStopRequest():

			return

		case ftpCmd := <-outChan:

			if ftpCmd == "test" {
				go func() {
					now := time.Now().String()
					path := "/chuis-le-goat" + now
					inChan <- view.AddNewProgressBar(path)

					var percent float64 = 0

					for {
						time.Sleep(40 * time.Millisecond)
						inChan <- progressbar.SetProgressFor(path, percent)
						percent += 0.002
					}
				}()
				continue
			}
			

			inChan <- log.AddLog(log.Usr, ftpCmd)
			_, err := conn.Write([]byte(ftpCmd + "\n"))
			if err != nil {
				ioErrorChan <- err
				return
			}

			if ftpCmd == "End" {
				connEndChan <- struct{}{}
			}
		}
	}
}
