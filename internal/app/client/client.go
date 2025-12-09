package client

import (
	"fmt"
	"ftp/internal/app/client/view"
	"ftp/internal/app/client/view/progressbar"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func Run(remote string) {

	outChan := make(chan string, 10)
	inChan := make(chan tea.Msg, 10)

	mainView := view.New(inChan, outChan)

	go func() {
		inChan <- view.AddNewProgressBar("/chuis-le-goat")

		var percent float64 = 0

		for {
			time.Sleep(40 * time.Millisecond)
			inChan <- progressbar.SetProgressFor("/chuis-le-goat", percent)
			percent += 0.01
		}
	}()

	go func() {
		inChan <- view.AddNewProgressBar("/chuis-le-goat2")

		var percent float64 = 0

		for {
			time.Sleep(50 * time.Millisecond)
			inChan <- progressbar.SetProgressFor("/chuis-le-goat2", percent)
			percent += 0.01
		}
	}()

	go func() {
		inChan <- view.AddNewProgressBar("/chuis-le-goat3")

		var percent float64 = 0

		for {
			time.Sleep(60 * time.Millisecond)
			inChan <- progressbar.SetProgressFor("/chuis-le-goat3", percent)
			percent += 0.01
		}
	}()

	if _, err := tea.NewProgram(mainView).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
	
}
