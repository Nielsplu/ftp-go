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
    Type    PacketType
    Data    any
}

type PacketOut struct {
	Buffer 	 []byte
	Priority bool
}