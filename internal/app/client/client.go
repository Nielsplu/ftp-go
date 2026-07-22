//go:build !js

package client

import (
	"ftp/internal/app/client/app"
	"ftp/internal/app/client/view"
	"ftp/internal/app/client/view/log"
	"log/slog"
	"net"
	
	tea "github.com/charmbracelet/bubbletea"
)

func Run(remote string) {

	conn, err := net.Dial("tcp", remote)
	if err != nil {
		slog.Error("Connection failed")
		return 
	}

	outChan := make(chan string, 10)
	inChan := make(chan tea.Msg, 10)

	mainView := view.New(inChan, outChan)

	// notify IU that connection is successful
	inChan <- log.AddLog(log.Sys, "Connecté à " + remote)

	go app.Handle(outChan, inChan, conn)
	tea.NewProgram(mainView).Run()

}
