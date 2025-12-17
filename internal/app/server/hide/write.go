package hide

import (
	"log/slog"
	"os"
	"fmt"
)

func writeTo(path string, hiddenFilesMap map[string]struct{}) {
	
    file, err := os.OpenFile(path, os.O_CREATE | os.O_WRONLY, 0600)
    if err != nil {
        slog.Error(err.Error())
    }

	defer file.Close()

	for filepath, _ := range hiddenFilesMap {
		_, err = file.WriteString(fmt.Sprintf("%s\n", filepath))
		if err != nil {
			slog.Error(err.Error())
		}
	}
}