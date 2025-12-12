package packet

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type fileChecked struct {
	clientPath, absPath string
	isDir				bool
}

func checkPath(rootPath, path string) (f fileChecked, err error) {

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

	// check if trying to escape root dir
	if !strings.HasPrefix(cleanPath, absRoot) {
		slog.Warn("error path traversal attempt", "path", path)
		return f, errors.New("trying to escape root dir")
	}

	// check if file exist
	fileSate, err := os.Stat(cleanPath)
	if err != nil {
		return f, errors.New("file doesn't exist")
	}

	return fileChecked{
		clientPath: cleanPath[len(absRoot):],
		absPath: cleanPath,
		isDir: fileSate.IsDir(),
	}, nil
}
