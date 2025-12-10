package hide

import (
	"log/slog"
	"os"
	"bufio"
)

func readFrom(path string) map[string]struct{} {

	file, err := os.Open(path)
	if err != nil {
        slog.Error(err.Error())
		return make(map[string]struct{})
    }

    defer file.Close()

    scanner := bufio.NewScanner(file)

	arrayOfHiddenFile := make(map[string]struct{})

    for scanner.Scan() {
        line := scanner.Text()
		arrayOfHiddenFile[line] = struct{}{}
    }

    if err := scanner.Err(); err != nil {
        slog.Error(err.Error())
		return make(map[string]struct{})
    }

	return arrayOfHiddenFile
	
}