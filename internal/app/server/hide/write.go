package hide

import (
	"fmt"
	"log/slog"
	"os"
)

func writeTo(path string, hiddenFilesMap map[string]struct{}) {
	// O_TRUNC : sans lui, réécrire une liste plus courte (après un Reveal)
	// laissait en place les octets de l'ancienne, et le fichier « démasqué »
	// réapparaissait masqué à la session suivante.
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		// Sans ce retour, le defer Close déréférençait un file nil.
		slog.Error(err.Error())
		return
	}

	defer file.Close()

	for filepath := range hiddenFilesMap {
		if _, err := file.WriteString(fmt.Sprintf("%s\n", filepath)); err != nil {
			slog.Error(err.Error())
		}
	}
}
