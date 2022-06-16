package lookergo

const workspacesSetBasePath = "4.0/workspaces"

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
