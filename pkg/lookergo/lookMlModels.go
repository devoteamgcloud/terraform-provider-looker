package lookergo

import "context"

const lookMlModelsBasePath = "4.0/lookml_models"

type LookMlModelsResource interface {
	List(ctx context.Context) ([]LookMLModel, *Response, error)
	Get(ctx context.Context, LookMLModelName string) (*LookMLModel, *Response, error)
	Create(ctx context.Context, LookMLModel *LookMLModel) (*LookMLModel, *Response, error)
	Update(ctx context.Context, LookMLModelName string, LookMLModel *LookMLModel) (*LookMLModel, *Response, error)
	Delete(ctx context.Context, LookMLModelName string) (*Response, error)
}

type LookmlModelNavExplore struct {
	Name        *string `json:"name,omitempty"`         // Name of the explore
	Description *string `json:"description,omitempty"`  // Description for the explore
	Label       *string `json:"label,omitempty"`        // Label for the explore
	Hidden      *bool   `json:"hidden,omitempty"`       // Is this explore marked as hidden
	GroupLabel  *string `json:"group_label,omitempty"`  // Label used to group explores in the navigation menus
  }

type LookMLModel struct {
	Can                      map[string]bool         `json:"can,omitempty"`                          // Operations the current user is able to perform on this object
	AllowedDbConnectionNames []string                `json:"allowed_db_connection_names,omitempty"`  // Array of names of connections this model is allowed to use
	Explores                 *[]LookmlModelNavExplore `json:"explores,omitempty"`                     // Array of explores (if has_content)
	HasContent               bool                    `json:"has_content,omitempty"`                  // Does this model declaration have have lookml content?
	Label                    string                  `json:"label,omitempty"`                        // UI-friendly name for this model
	Name                     string                  `json:"name,omitempty"`                         // Name of the model. Also used as the unique identifier
	ProjectName              string                  `json:"project_name,omitempty"`                 // Name of project containing the model
	UnlimitedDbConnections   bool                    `json:"unlimited_db_connections,omitempty"`     // Is this model allowed to use all current and future connections
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
