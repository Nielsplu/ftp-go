package server

import (
	"fmt"
	"ftp/internal/app/server/packet"
<<<<<<< HEAD
	"ftp/internal/app/server/remote"
	t "ftp/internal/app/server/types"
	"ftp/internal/pkg/utils"
	"log/slog"
=======
	"ftp/internal/pkg/utils"
	"ftp/internal/app/server/hide"

	t "ftp/internal/app/server/types"
>>>>>>> 8951daf (update Hide Reveal and List)
)

func Start(config t.ServerConfig) {

	mainStopper := utils.NewStopper()
	packetInChan := make(chan t.PacketIn, 10)
	internalConnActionChan := make(chan t.InternalConnAction, 10)

	// listen on non admin port
	mainStopper.Go(func(child *utils.Stopper) {
		remote.ListenOn(&config.Port, false, packetInChan, internalConnActionChan, child)
	})

	// listen on admin port
	mainStopper.Go(func(child *utils.Stopper) {
		remote.ListenOn(&config.AdminPort, true, packetInChan, internalConnActionChan, child)
	})

	config.HiddenFiles = hide.ReadFromFile("hiddenFile.txt")
	state := t.ServerState {
		Clients: make(map[string]t.Client),
		Config: config,
	}

	for {

		select {
		case connAction := <- internalConnActionChan:

			switch connAction := connAction.(type) {
			case t.NewConn:

				state.Clients[connAction.Id] = connAction.Client
				slog.Info(fmt.Sprintf("%d clients connected", len(state.Clients)))

			case t.ConnEnd:

				delete(state.Clients, connAction.Id)
				slog.Info(fmt.Sprintf("%d clients connected", len(state.Clients)))

			}

		case packetIn := <-packetInChan:

			client, exists := state.Clients[packetIn.ClientId]
			if !exists {
				slog.Error("No client with id : " + packetIn.ClientId)
				continue
			}

			switch packetIn.Type {

			case t.List:
				
				packet.PerformeList(packetIn.Path, client.PacketOutChan, config.RootPath, config.HiddenFiles)

			case t.Get:

				packet.PerformeGet(packetIn.Path, client.PacketOutChan, config.RootPath, mainStopper)

			case t.Hide:

				packet.PerformeHide(packetIn.Path, client.PacketOutChan, config.RootPath, config.HiddenFiles)

			case t.Reveal:

				packet.PerformeReveal(packetIn.Path, client.PacketOutChan, config.RootPath, config.HiddenFiles)

			case t.Cd:

				packet.PerformeCd(packetIn.Path, client.PacketOutChan, config.RootPath, &client)

			case t.Terminate:

				mainStopper.StopChilds()
				println("Bye bye !")
				return

			default:
				slog.Error(fmt.Sprintf("Unknown packetType : %d", packetIn.Type))
			}

		}
	}

}
