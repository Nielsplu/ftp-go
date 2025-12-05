package server

import (
	"fmt"
	"io"
	"net"
	"os"
)

// HandleGet s'occupe d'envoyer un fichier au client
func HandleGet(conn net.Conn, nomFichier string) {
	
	fichier, err := os.Open(nomFichier)
	if err != nil {
		fmt.Fprintln(conn, "FileUnknown")
		return 
	}
	//ferme le fichier à la fin de la fonction
	defer fichier.Close()

	fmt.Fprint(conn, "Start\n")
	
	//tableau d'octets temporaire.
	//on déplace 1024 octets.
    buffer := make([]byte, 1024)

	for {
		// n le nombre d'octets 
		n, err := fichier.Read(buffer)

		if err == io.EOF{
			break
		}
		if err != nil{
			fmt.Println("Erreur de lecture du fichier:", err)	
		}

        // buffer[:n] pour le cas ou le dernier
		//  morceau est plus petit que 1024 octetss.
        _, err = conn.Write(buffer[:n])

		if err != nil {
            fmt.Println("Erreur d'envoi sur le réseau:", err)
            return
        }

	}
}