package view

// Accesseurs utilisés par le pont WebAssembly (cmd/wasm).

func (m NewProgressBarMsg) Filename() string { return m.filename }

func (m RemoveProgressBarMsg) Filename() string { return m.filename }

func (e FatalError) Message() string { return e.message }

func (e FatalError) Details() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}
