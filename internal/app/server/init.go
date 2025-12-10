package server

import (
	"fmt"
	"log/slog"
	"ftp/internal/app/server/remote"
	"ftp/internal/app/server/packet"
	t "ftp/internal/app/server/types"
	"ftp/internal/pkg/utils"
)

func Start(config ServerConfig) {

	mainStopper := utils.NewStopper()
	packetInChan := make(chan t.PacketIn, 10)

	// listen on non admin port
	mainStopper.Go(func(child *utils.Stopper) { 
		remote.ListenOn(&config.Port, false, packetInChan, child)
	})

	// listen on admin port
	mainStopper.Go(func(child *utils.Stopper) { 
		remote.ListenOn(&config.AdminPort, true, packetInChan, child)
	})

	nClientConnected := 0

	for {
		packetIn := <- packetInChan

		switch packetIn.Type {
		case t.List: packet.PerformeList(packetIn.Path, packetIn.AnswerChan)
		case t.Get: packet.PerformeGet(packetIn.Path, packetIn.AnswerChan)
		case t.Hide: packet.PerformeHide(packetIn.Path, packetIn.AnswerChan)
		case t.Reveal: packet.PerformeReveal(packetIn.Path, packetIn.AnswerChan)
		case t.NewConn: 

			nClientConnected += 1
			slog.Info(fmt.Sprintf("%d clients connected", nClientConnected))

		case t.ConnEnd: 


			nClientConnected -= 1
			slog.Info(fmt.Sprintf("%d clients connected", nClientConnected))

		case t.Terminate:
			
			mainStopper.StopChilds()
			println("Bye bye !")
			return

		default:
			slog.Error(fmt.Sprintf("Unknown packetType : %d", packetIn.Type))
		}
	}

}