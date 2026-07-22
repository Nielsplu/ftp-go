//go:build js && wasm

package wasmbridge

import (
	"net"
	"syscall/js"
	"time"

	"ftp/internal/app/client/app"
	"ftp/internal/app/client/view"
	"ftp/internal/app/client/view/cmd"
	"ftp/internal/app/client/view/log"
	"ftp/internal/app/client/view/progressbar"
	"ftp/internal/app/server"
	t "ftp/internal/app/server/types"

	tea "github.com/charmbracelet/bubbletea"
)

// Session en cours : le writer du client lit sur ce channel.
var outChan chan string

func demarrerServeur() {
	go server.Start(t.ServerConfig{Port: "3333", AdminPort: "4444", RootPath: "data/"})
}

// connect(rappel[, port]) : ouvre une session client (port 3333 par défaut,
// 4444 pour l'admin, comme le flag -p du client natif) et relaie chaque
// message de la vue vers rappel(type, ...args). Relance le serveur s'il a été
// arrêté par un Terminate précédent.
func connecter(_ js.Value, args []js.Value) any {
	rappel := args[0]
	port := "3333"
	if len(args) > 1 {
		port = args[1].String()
	}
	adresse := "127.0.0.1:" + port

	go func() {
		var conn net.Conn
		var err error
		for essai := 0; essai < 20; essai++ {
			conn, err = net.Dial("tcp", adresse)
			if err == nil {
				break
			}
			if essai == 10 {
				demarrerServeur()
			}
			time.Sleep(100 * time.Millisecond)
		}
		if err != nil {
			rappel.Invoke("fatal", "Connexion au serveur impossible", err.Error())
			return
		}

		inChan := make(chan tea.Msg, 10)
		out := make(chan string, 10)
		outChan = out

		go app.Handle(out, inChan, conn)
		rappel.Invoke("log", "Sys", "Connecté à "+adresse)

		for msg := range inChan {
			switch m := msg.(type) {
			case log.LogMsg:
				rappel.Invoke("log", m.OriginLabel(), m.Content())
			case cmd.ChangeDirMsg:
				rappel.Invoke("cwd", m.To)
			case view.NewProgressBarMsg:
				rappel.Invoke("progress", m.Filename(), 0.0)
			case progressbar.ProgressMsg:
				rappel.Invoke("progress", m.Filename, m.Percent)
			case view.RemoveProgressBarMsg:
				rappel.Invoke("progressEnd", m.Filename())
			case view.FatalError:
				rappel.Invoke("fatal", m.Message(), m.Details())
				return
			case tea.QuitMsg:
				rappel.Invoke("quit")
				return
			}
		}
	}()

	return nil
}

// send(ligne) : transmet une commande utilisateur au writer du client.
func envoyer(_ js.Value, args []js.Value) any {
	if outChan != nil {
		outChan <- args[0].String()
	}
	return nil
}

// Start expose le pont sous globalThis.__ftpgo puis bloque pour garder le
// runtime Go vivant.
func Start() {
	demarrerServeur()

	js.Global().Set("__ftpgo", js.ValueOf(map[string]any{}))
	objet := js.Global().Get("__ftpgo")
	objet.Set("connect", js.FuncOf(connecter))
	objet.Set("send", js.FuncOf(envoyer))
	objet.Set("ready", true)

	select {}
}
