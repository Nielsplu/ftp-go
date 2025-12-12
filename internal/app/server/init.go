package server

import (
	"fmt"
	"ftp/internal/app/server/packet"
	"ftp/internal/app/server/remote"
	"ftp/internal/app/server/hide"
	t "ftp/internal/app/server/types"
	"ftp/internal/pkg/utils"
	"log/slog"
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

	state := t.ServerState {
		Clients: make(map[string]*t.Client),
		Config: config,
		HiddenFiles: hide.ReadFromFile("hiddenFile.txt"),
	}

	for {

		select {
		case connAction := <- internalConnActionChan:

			switch connAction := connAction.(type) {
			case t.NewConn:

				state.Clients[connAction.Id] = &connAction.Client
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
				
				packet.PerformeList(packetIn.Path, client.PacketOutChan, &state, client)

			case t.Get:

				packet.PerformeGet(packetIn.Path, client.PacketOutChan, &state, client, mainStopper)

			case t.Hide:

				packet.PerformeHide(packetIn.Path, client.PacketOutChan, &state, client)

			case t.Reveal:

				packet.PerformeReveal(packetIn.Path, client.PacketOutChan, &state, client)

			case t.Cd:

				packet.PerformeCd(packetIn.Path, client.PacketOutChan, &state, client)

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
