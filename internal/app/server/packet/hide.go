package packet

import (
	t "ftp/internal/app/server/types"
	"log/slog"
	"os"
	"fmt"
	"bufio"
)

func PerformeHide(path string, responseChan chan t.PacketOut) {
	
}

func CreateArrayOfHidenFile(/*file string*/) /*[]string*/{
	dir, err := os.Getwd()
	if err != nil {
		slog.Error(err.Error())
	}

	dir += "/../../internal/app/server/"

	file, err := os.Open(dir + "hidden_file.txt")
	if err != nil {
        slog.Error(err.Error())
		return
    }

    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(line)
    }

	fmt.Println("here")
    if err := scanner.Err(); err != nil {
        slog.Error(err.Error())
		return
    }
	
}