package log

// Accesseurs utilisés par le pont WebAssembly (cmd/wasm) : la vue web doit
// pouvoir lire les messages sans dépendre du rendu Bubbletea.

func (m LogMsg) Content() string { return m.content }

func (m LogMsg) OriginLabel() string {
	switch m.origin {
	case Usr:
		return "Usr"
	case Sys:
		return "Sys"
	case Srv:
		return "Srv"
	}
	return "?"
}
