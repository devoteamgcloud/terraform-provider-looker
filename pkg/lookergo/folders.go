package lookergo

const foldersSetBasePath = "api/4.0/folders"

type FoldersResource interface {
}

type FoldersResourceOp struct {
	client *Client
}

var _ FoldersResource = &FoldersResourceOp{}

/*
List
ListById
ListByName
Get -> recurse children.
Create
Update
Delete
*/
