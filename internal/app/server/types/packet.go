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

type Client struct {
	CurrentPath   string
	PacketOutChan chan PacketOut
}

type ServerConfig struct {
	Port, AdminPort, RootPath string
}

type ServerState struct {
	Clients map[string]Client
	Config  ServerConfig
}

type PacketIn struct {
	Type           PacketType
	Path, ClientId string
}

type PacketOut struct {
	Buffer      []byte
	LowPriority bool
}
