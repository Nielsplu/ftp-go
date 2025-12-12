package types

type PacketType int

const (
	List PacketType = iota
	Get
	Cd
	End
	Hide
	Reveal
	Terminate
)

type PacketIn struct {
	Type           PacketType
	Path, ClientId string
}

type PacketOut struct {
	Buffer      []byte
	LowPriority bool
}
