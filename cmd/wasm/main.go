//go:build js && wasm

package main

import (
	bridge "ftp/wasm"
)

func main() {
	bridge.Start()
}
