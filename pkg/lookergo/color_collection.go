package lookergo
​
import "context"
​
const ColorCollectionBasePath = "4.0/color_collections"
// https://developers.looker.com/api/explorer/4.0/types/ColorCollection/ColorCollection
type ColorCollection struct {
	Id                  string              `json:"id,omitempty"`                  // Unique Id
	Label               string              `json:"label,omitempty"`               // Label of color collection
	CategoricalPalettes *[]DiscretePalette   `json:"categoricalPalettes,omitempty"` // Array of categorical palette definitions
	SequentialPalettes  *[]ContinuousPalette `json:"sequentialPalettes,omitempty"`  // Array of discrete palette definitions
	DivergingPalettes   *[]ContinuousPalette `json:"divergingPalettes,omitempty"`   // Array of diverging palette definitions
}
​
// https://developers.looker.com/api/explorer/4.0/types/ColorCollection/DiscretePalette
type DiscretePalette struct {
	Id     string   `json:"id,omitempty"`     // Unique identity string
	Label  string   `json:"label,omitempty"`  // Label for palette
	Type   string   `json:"type,omitempty"`   // Type of palette
	Colors []string `json:"colors,omitempty"` // Array of colors in the palette
}
​
// https://developers.looker.com/api/explorer/4.0/types/ColorCollection/ContinuousPalette
type ContinuousPalette struct {
	Id    string      `json:"id,omitempty"`    // Unique identity string
	Label string      `json:"label,omitempty"` // Label for palette
	Type  string      `json:"type,omitempty"`  // Type of palette
	Stops *[]ColorStop `json:"stops,omitempty"` // Array of ColorStops in the palette
}
​
// https://developers.looker.com/api/explorer/4.0/types/ColorCollection/WriteColorCollection
type WriteColorCollection struct {
	Label               string              `json:"label,omitempty"`               // Label of color collection
	CategoricalPalettes *[]DiscretePalette   `json:"categoricalPalettes,omitempty"` // Array of categorical palette definitions
	SequentialPalettes  *[]ContinuousPalette `json:"sequentialPalettes,omitempty"`  // Array of discrete palette definitions
	DivergingPalettes   *[]ContinuousPalette `json:"divergingPalettes,omitempty"`   // Array of diverging palette definitions
}
​
// https://developers.looker.com/api/explorer/4.0/types/ColorCollection/ColorStop
type ColorStop struct {
	Color  string `json:"color,omitempty"`  // CSS color string
	Offset int  `json:"offset,omitempty"` // Offset in continuous palette (0 to 100)
}
​
type ColorCollectionResourceOp struct {
	client *Client
}
​
​
type ColorCollectionResource interface {
	List(context.Context, *ListOptions) ([]ColorCollection, *Response, error)
	Get(context.Context, string) (*ColorCollection, *Response, error)
	Create(context.Context, *WriteColorCollection) (*ColorCollection, *Response, error)
	Update(context.Context, string, *ColorCollection) (*ColorCollection, *Response, error)
	Delete(context.Context, string) (*Response, error)
}
​
func (s *ColorCollectionResourceOp) List(ctx context.Context, opt *ListOptions) ([]ColorCollection, *Response, error) {
	return doList(ctx, s.client, ColorCollectionBasePath, opt, new([]ColorCollection))
}
​
func (s *ColorCollectionResourceOp) Get(ctx context.Context, ColorCollectionId string) (*ColorCollection, *Response, error) {
	return doGetById(ctx, s.client, ColorCollectionBasePath, ColorCollectionId, new(ColorCollection))
}
​
func (s *ColorCollectionResourceOp) Create(ctx context.Context, requestColorCollection *WriteColorCollection) (*ColorCollection, *Response, error) {
	return doCreate(ctx, s.client, ColorCollectionBasePath, requestColorCollection, new(ColorCollection))
}
​
func (s *ColorCollectionResourceOp) Update(ctx context.Context, ColorCollectionId string, requestColorCollection *ColorCollection) (*ColorCollection, *Response, error) {
	return doUpdate(ctx, s.client, ColorCollectionBasePath, ColorCollectionId, requestColorCollection, new(ColorCollection))
}
​
func (s *ColorCollectionResourceOp) Delete(ctx context.Context, ColorCollectionId string) (*Response, error) {
	return doDelete(ctx, s.client, ColorCollectionBasePath, ColorCollectionId)
}