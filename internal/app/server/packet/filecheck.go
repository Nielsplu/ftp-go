package packet

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	t "ftp/internal/app/server/types"
)

type fileChecked struct {
	clientVisiblePath, absPath string
	isDir, isHidden	    	   bool
}

func checkPath(state *t.ServerState, clientPath, path string) (f fileChecked, err error) {

	//get the absolute path
	absRoot, err := filepath.Abs(state.Config.RootPath)
	if err != nil {
		slog.Error("Invalid root path")
		return
	}

	//get the path
	var fullPath string
	if len(path) > 0 && path[0] == '/' {
		fullPath = filepath.Join(absRoot, path)
	}else {
		fullPath = filepath.Join(absRoot, clientPath, path)
	}

	//clean path
	cleanPath, err := filepath.Abs(filepath.Clean(fullPath))
	if err != nil {
		slog.Error("error getting abs path")
		return
	}

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

	clientVisiblePath := cleanPath[len(absRoot):]
	return fileChecked{
		clientVisiblePath: clientVisiblePath,
		absPath: cleanPath,
		isDir: fileSate.IsDir(),
		isHidden: state.HiddenFiles.IsPathHidden(clientVisiblePath),
	}, nil
}
