package types

type PacketType int

const (
	List PacketType = iota
	Get
	End
	Hide
	Reveal
	Terminate
	NewConn
	ConnEnd
)

type PacketIn struct {
	Type       PacketType
	Path       string
	AnswerChan chan PacketOut
}

type PacketOut struct {
	Buffer      []byte
	LowPriority bool
}
