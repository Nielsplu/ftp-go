package types

type NewConn struct {
	Client Client
}

type ConnEnd struct {
	Id string
}

type InternalConnAction interface{}
