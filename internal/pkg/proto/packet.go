package proto

type Packet int

const (
	List Packet = iota
	Get
	End
	Hide
	Reveal
	Terminate
)
