// Stub minimal de bubbletea pour le build WebAssembly : la vraie TUI ne
// compile pas en js/wasm (TTY, presse-papier), mais la logique client
// n'utilise que ces quatre symboles. Le build natif garde la vraie
// dépendance ; ce module n'est branché que par le replace de cmd/wasm/go.mod.
package tea

type Msg interface{}

type Cmd func() Msg

type QuitMsg struct{}

func Quit() Msg { return QuitMsg{} }
