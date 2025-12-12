package remote

import (
	"bufio"
	"ftp/internal/app/server/packet"
	t "ftp/internal/app/server/types"
	"ftp/internal/pkg/utils"
	"log/slog"
	"net"
	"time"
)

func Handle(
	conn net.Conn,
	admin bool,
	internalConnActionChan chan t.InternalConnAction,
	packetInChan chan t.PacketIn,
	stopper *utils.Stopper,
) bool {
	

	id := utils.MakeId()
	packetOutChan := make(chan t.PacketOut, 10)

	resetTimerChan := make(chan struct{}, 10)

	// notify app that new conn as started
	internalConnActionChan <- t.NewConn { 
		Id: id,
		Client: t.Client{
			CurrentPath: "/",
			PacketOutChan: packetOutChan,
		},
	}

	// Do not defer conn.Close here
	// it's writer responsability to 
	// close the connection, like 
	// it needs to send END\n before

	// channel used to coordinate writer 
	// and reader when an error occurs
	// set capacity to one to avoid blocking
	ioErrorChan := make(chan error, 1)

	// create writer goroutine with priority
	stopper.Go(func(child *utils.Stopper) {
		
		err := StartWriter(conn, packetOutChan, packetInChan, child, resetTimerChan)

		// send error if reader 
		// hasn't send one yet
		if len(ioErrorChan) < cap(ioErrorChan) {
			ioErrorChan <- err
		}

	})

	reader := bufio.NewReader(conn)
	lineChan := make(chan string, 1)

	timeoutDuration := 5 * time.Second
    timer := time.NewTimer(timeoutDuration)

	defer timer.Stop()

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

		case <-timer.C:
            // time whithout activity
            stopper.StopChilds()
            slog.Debug("Conn shutdown (due to inactivity)")
            
            // say connection end
            internalConnActionChan <- t.ConnEnd { Id: id }
            return false

		case <- resetTimerChan:

			// stop the timer.
			if !timer.Stop() {

				// try receive timer signal
				// if it hs been send while
				// writing line (not blocking)
                select {
                case <-timer.C:
                default:
                }
            }

			timer.Reset(timeoutDuration)

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
			internalConnActionChan <- t.ConnEnd { Id: id }
			return false

		case line := <-lineChan:

			// reset timer.
			if !timer.Stop() {
				// receive (not blocking) timer signal
				// like it has been send while reading line
                select {
                case <-timer.C:
                default:
                }
            }

			timer.Reset(timeoutDuration)

			// ignore Ok
			if line == "OK\n" {
				continue
			}

			// parse cmd
			packetIn, err := packet.Parse(line, admin)
			if err != nil {
				packetOutChan <- t.PacketOut{ Buffer: []byte(err.Error() + "\n") }
				slog.Debug("Parse error : " + err.Error())
				continue
			}

			if packetIn.Type == t.End {

				// stop writer and reader
				stopper.StopChilds()
				slog.Debug("Conn shutdown (End requested)")

				// notify core that conn has ended
				internalConnActionChan <- t.ConnEnd { Id: id }
				return false
			}

			packetIn.ClientId = id
			packetInChan <- packetIn
		}
	}
}
