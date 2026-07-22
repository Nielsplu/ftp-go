// Package wasmbridge fait tourner le vrai serveur et le vrai client FTP dans
// un navigateur, reliés par le réseau loopback en mémoire de Go (paquet net
// sur js/wasm). Seule la vue Bubbletea est remplacée par des callbacks
// JavaScript : reader, writer, stoppers et protocole sont strictement ceux
// des binaires natifs.
//
// Compilation depuis cmd/wasm (module dédié qui remplace bubbletea par un
// stub, la TUI ne compilant pas en js/wasm) :
//
//	cd cmd/wasm && GOOS=js GOARCH=wasm go build -o ftp.wasm .
//
// Le système de fichiers (data/, hiddenFile.txt, downloads/) est fourni par
// la page hôte via le shim globalThis.fs de wasm_exec.js.
package wasmbridge
