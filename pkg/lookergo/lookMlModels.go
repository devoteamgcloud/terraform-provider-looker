package lookergo

import "context"

const lookMlModelsBasePath = "4.0/models"

type LookMlModelsResource interface {
	List(ctx context.Context) ([]ModelSet, *Response, error)
	Get(ctx context.Context, modelSetId string) (*ModelSet, *Response, error)
	Create(ctx context.Context, modelSet *ModelSet) (*ModelSet, *Response, error)
	Update(ctx context.Context, modelSetId string, modelSet *ModelSet) (*ModelSet, *Response, error)
	Delete(ctx context.Context, modelSetId string) (*Response, error)
}

func (s LookMlModelsResourceOp) List(ctx context.Context) ([]ModelSet, *Response, error) {
	// TODO implement me
	panic("implement me")
}

func (s LookMlModelsResourceOp) Get(ctx context.Context, modelSetId string) (*ModelSet, *Response, error) {
	// TODO implement me
	panic("implement me")
}

func (s LookMlModelsResourceOp) Create(ctx context.Context, modelSet *ModelSet) (*ModelSet, *Response, error) {
	// TODO implement me
	panic("implement me")
}

func (s LookMlModelsResourceOp) Update(ctx context.Context, modelSetId string, modelSet *ModelSet) (*ModelSet, *Response, error) {
	// TODO implement me
	panic("implement me")
}

func (s LookMlModelsResourceOp) Delete(ctx context.Context, modelSetId string) (*Response, error) {
	// TODO implement me
	panic("implement me")
}

type LookMlModelsResourceOp struct {
	client *Client
}

var _ LookMlModelsResource = &LookMlModelsResourceOp{}
