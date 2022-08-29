package lookergo

import "context"

const lookMlModelsBasePath = "api/4.0/lookml_models"


type LookMlModelsResource interface {
	List(ctx context.Context) ([]LookMLModel, *Response, error)
	Get(ctx context.Context, LookMLModelName string) (*LookMLModel, *Response, error)
	Create(ctx context.Context, LookMLModel *LookMLModel) (*LookMLModel, *Response, error)
	Update(ctx context.Context, LookMLModelName string, LookMLModel *LookMLModel) (*LookMLModel, *Response, error)
	Delete(ctx context.Context, LookMLModelName string) (*Response, error)
}

type LookMLModel struct {
	Name                        string   `json:"name,omitempty"`
	Project_name                string   `json:"project_name,omitempty"`
	Label                       string   `json:"label,omitempty"`
	Allowed_db_connection_names []string `json:"allowed_db_connection_names,omitempty"`
	Unlimited_db_connections    bool     `json:"unlimited_db_connections,omitempty" `
}

func (s LookMlModelsResourceOp) List(ctx context.Context) ([]LookMLModel, *Response, error) {
	// TODO implement me
	//return doList(ctx, s.client, lookMlModelsBasePath, )
	panic("Not implemented")
}

func (s LookMlModelsResourceOp) Get(ctx context.Context, LookMLModelName string) (*LookMLModel, *Response, error) {
	return doGetById(ctx, s.client, lookMlModelsBasePath, LookMLModelName, new(LookMLModel))
}

func (s LookMlModelsResourceOp) Create(ctx context.Context, requestLookMLModel *LookMLModel) (*LookMLModel, *Response, error) {
	return doCreate(ctx, s.client, lookMlModelsBasePath, requestLookMLModel, new(LookMLModel))
}

func (s LookMlModelsResourceOp) Update(ctx context.Context, LookMLModelName string, requestLookMLModel *LookMLModel) (*LookMLModel, *Response, error) {
	return doUpdate(ctx, s.client, lookMlModelsBasePath, LookMLModelName, requestLookMLModel, new(LookMLModel))
}

func (s LookMlModelsResourceOp) Delete(ctx context.Context, LookMLModelName string) (*Response, error) {
	return doDelete(ctx, s.client, lookMlModelsBasePath, LookMLModelName)
}

type LookMlModelsResourceOp struct {
	client *Client
}

var _ LookMlModelsResource = &LookMlModelsResourceOp{}
