package types

type NewConn struct {
	Id     string
	Client Client
}

type ConnEnd struct {
	Id string
}

type InternalConnAction interface{}
