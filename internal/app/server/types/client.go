package types

type Client struct {
	Id, CurrentPath string
	PacketOutChan   chan PacketOut
}