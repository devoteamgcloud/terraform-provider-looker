package lookergo

// https://developers.looker.com/api/explorer/4.0/types/ColorCollection/ColorCollection
type ColorCollection struct {
	Id                  *string              `json:"id,omitempty"`                  // Unique Id
	Label               *string              `json:"label,omitempty"`               // Label of color collection
	CategoricalPalettes *[]DiscretePalette   `json:"categoricalPalettes,omitempty"` // Array of categorical palette definitions
	SequentialPalettes  *[]ContinuousPalette `json:"sequentialPalettes,omitempty"`  // Array of discrete palette definitions
	DivergingPalettes   *[]ContinuousPalette `json:"divergingPalettes,omitempty"`   // Array of diverging palette definitions
}

// https://developers.looker.com/api/explorer/4.0/types/ColorCollection/DiscretePalette
type DiscretePalette struct {
	Id     *string   `json:"id,omitempty"`     // Unique identity string
	Label  *string   `json:"label,omitempty"`  // Label for palette
	Type   *string   `json:"type,omitempty"`   // Type of palette
	Colors *[]string `json:"colors,omitempty"` // Array of colors in the palette
}

// https://developers.looker.com/api/explorer/4.0/types/ColorCollection/ContinuousPalette
type ContinuousPalette struct {
	Id    *string      `json:"id,omitempty"`    // Unique identity string
	Label *string      `json:"label,omitempty"` // Label for palette
	Type  *string      `json:"type,omitempty"`  // Type of palette
	Stops *[]ColorStop `json:"stops,omitempty"` // Array of ColorStops in the palette
}

// https://developers.looker.com/api/explorer/4.0/types/ColorCollection/WriteColorCollection
type WriteColorCollection struct {
	Label               *string              `json:"label,omitempty"`               // Label of color collection
	CategoricalPalettes *[]DiscretePalette   `json:"categoricalPalettes,omitempty"` // Array of categorical palette definitions
	SequentialPalettes  *[]ContinuousPalette `json:"sequentialPalettes,omitempty"`  // Array of discrete palette definitions
	DivergingPalettes   *[]ContinuousPalette `json:"divergingPalettes,omitempty"`   // Array of diverging palette definitions
}

// https://developers.looker.com/api/explorer/4.0/types/ColorCollection/ColorStop
type ColorStop struct {
	Color  *string `json:"color,omitempty"`  // CSS color string
	Offset *int64  `json:"offset,omitempty"` // Offset in continuous palette (0 to 100)
}

type ColorCollectionOp struct {
	client *Client
}

const ColorCollectionBasePath = "4.0/color_collection"
