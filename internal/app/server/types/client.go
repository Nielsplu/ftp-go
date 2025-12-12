package types

type Client struct {
	CurrentPath   string
	PacketOutChan chan PacketOut
}