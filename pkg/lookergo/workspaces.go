package lookergo

type WorkspacesResource interface {
}

type WorkspacesResourceOp struct {
	client *Client
}

var _ WorkspacesResource = &WorkspacesResourceOp{}

/*
List
Get
*/
