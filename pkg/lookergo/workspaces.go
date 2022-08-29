package lookergo

const workspacesSetBasePath = "api/4.0/workspaces"

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
