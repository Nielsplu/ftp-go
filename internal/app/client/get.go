package client

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
)

func processGet(
	conn net.Conn,
	reader *bufio.Reader,
	fileName string,
	fileSize int64) {

	message, _ := reader.ReadString('\n')
	message = strings.TrimSpace(message)

	if message == "FileUnknown" {
		fmt.Println("Erreur : Le serveur ne trouve pas ce fichier.")
		return
	}

	if message != "Start" {
		fmt.Println("Erreur du serveur :", message)
		return
	}

	LocalFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Impossible de créer le fichier:", err)
		return
	}
	defer LocalFile.Close()

	fmt.Printf("Je commence à télécharger %s (%d octets)...\n", fileName, fileSize)

	

}

