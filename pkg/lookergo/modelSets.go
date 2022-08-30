package lookergo

import "context"

const modelSetsBasePath = "4.0/model_sets"

type ModelSetsResource interface {
	List(ctx context.Context) ([]ModelSet, *Response, error)
	Get(ctx context.Context, modelSetId string) (*ModelSet, *Response, error)
	Create(ctx context.Context, modelSet *ModelSet) (*ModelSet, *Response, error)
	Update(ctx context.Context, modelSetId string, modelSet *ModelSet) (*ModelSet, *Response, error)
	Delete(ctx context.Context, modelSetId string) (*Response, error)
}

type ModelSetsResourceOp struct {
	client *Client
}

var _ ModelSetsResource = &ModelSetsResourceOp{}

type ModelSet struct {
	BuiltIn   bool     `json:"built_in,omitempty"`
	Id        string   `json:"id,omitempty"`
	AllAccess bool     `json:"all_access,omitempty"`
	Models    []string `json:"models,omitempty"`
	Name      string   `json:"name,omitempty"`
}

func (s ModelSetsResourceOp) List(ctx context.Context) ([]ModelSet, *Response, error) {
	return doList(ctx, s.client, modelSetsBasePath, nil, new([]ModelSet))
}

func (s ModelSetsResourceOp) Get(ctx context.Context, modelSetId string) (*ModelSet, *Response, error) {
	return doGetById(ctx, s.client, modelSetsBasePath, modelSetId, new(ModelSet))
}

func (s ModelSetsResourceOp) Create(ctx context.Context, modelSet *ModelSet) (*ModelSet, *Response, error) {
	return doCreate(ctx, s.client, modelSetsBasePath, modelSet, new(ModelSet))
}

func (s ModelSetsResourceOp) Update(ctx context.Context, modelSetId string, modelSet *ModelSet) (*ModelSet, *Response, error) {
	return doUpdate(ctx, s.client, modelSetsBasePath, modelSetId, modelSet, new(ModelSet))
}

func (s ModelSetsResourceOp) Delete(ctx context.Context, modelSetId string) (*Response, error) {
	return doDelete(ctx, s.client, modelSetsBasePath, modelSetId)
}
