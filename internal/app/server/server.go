package server

import (
	"log/slog"
	"net"
	"fmt"
	"os"
	"bufio"
	"strings"
)
var nbClient=0
func StartServer(port *string, stopper *Stopper) {

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("Now listening on port " + *port)

	connChan := make(chan net.Conn)
	errorChan := make(chan error)

	for {

		// start a stopper goroutine to listen, 
		// like that when shutdwon is required 
		// listener will send error in chanel
		// and StartServer will end
		stopper.Go(func(_ *Stopper) {
			conn, err := listener.Accept()
			if err != nil {
				errorChan <- err
				return
			}
			connChan <- conn
		})

		select {
		case <-stopper.Wait():

			slog.Debug("Stopping accept loop ...")
			listener.Close()
			stopper.Stop()
			return

		case conn := <- connChan:

			stopper.Go(func (child *Stopper){ 
				Handle(conn, child) 
			})
		
		case err := <- errorChan:
			slog.Error(err.Error())
		}
	}
}




func ReturnFile(conn  net.Conn , nomFichier string){
	fichier, err := os.Open(nomFichier)
	if err != nil{
		fmt.Println(conn, "FileUnknown")
		return 
	}

	defer fichier.Close()
	fmt.Print(conn, "Start\n")
}

func GererClient(conn net.Conn){
	nbClient++
	slog.Info("Client connecté","Total", nbClient)

	defer func ()  {
		nbClient--
		slog.Info("Client parti","Total", nbClient)
		conn.Close()
	}()

	scanner := bufio.NewScanner(conn)


	for scanner.Scan(){
		ligne:=scanner.Text()
		mots := strings.Fields(ligne)

		if len(mots)==0{
			continue
		}

		commande:= mots[0]

		switch commande{
			case "Get":
				if len(mots) < 2 {
					fmt.Fprintln(conn, "Erreur Nom fichier manquant")
					continue
				}
				ReturnFile(conn,mots[0])
			}
	}
}