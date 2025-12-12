package packet

import (
	"errors"
	"log/slog"
	"path/filepath"
	"strings"
)

func checkPath(rootPath, path string) (s string, err error) {
	
	//get the absolute path
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		slog.Error("Invalid root path")
		return
	}

	//get the path
	fullPath := filepath.Join(absRoot, path)

	//clean path
	cleanPath := filepath.Clean(fullPath)

	if !strings.HasPrefix(cleanPath, absRoot) {
		slog.Warn("error path traversal attempt", "path", path)
		return s, errors.New("trying to escape root dir")
	}

	return cleanPath[len(absRoot):], nil
}