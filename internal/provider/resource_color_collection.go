package provider

import (
	"context"
	//"fmt"
	//"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	//"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	//"net/http"
)

func resourceColorCollection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceColorCollectionCreate,
		ReadContext:   resourceColorCollectionRead,
		UpdateContext: resourceColorCollectionUpdate,
		DeleteContext: resourceColorCollectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ColorCollection id",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"label": {
				Description: "Label of color collection",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"categoricalPalettes": {
				Description: "Array of categorical palette definitions",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Label of color collection",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"label": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"colors": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"sequentialPalettes": {
				Description: "Array of categorical palette definitions",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Unique ID of palette",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"label": {
							Type:        schema.TypeString,
							Description: "Label of palette",
							Optional:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "type of palette",
						},
						"stops": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Array of ColorStops in the palette",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"color": {
										Type:        schema.TypeString,
										Description: "CSS color string",
										Required:    true,
									},
									"offset": {
										Type:         schema.TypeInt,
										Description:  "Offset in continuous palette (0 to 100)",
										Default:      0,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 100),
									},
								},
							},
						},
					},
				},
			},
			"divergingPalettes": {
				Description: "Array of categorical palette definitions",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Unique identity string",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"label": {
							Description: "Label for palette",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Type of palette",
							Optional:    true,
						},
						"stops": {
							Type:        schema.TypeList,
							Description: "Array of ColorStops in the palette",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"color": {
										Type:        schema.TypeString,
										Description: "CSS color string",
										Required:    true,
									},
									"offset": {
										Type:         schema.TypeInt,
										Description:  "Offset in continuous palette (0 to 100)",
										Default:      0,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 100),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Receives terraform resource schema, builds a golang struct with json fields from it, sends a Post request with the
func resourceColorCollectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	return nil
}

func resourceColorCollectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	return nil
}

func resourceColorCollectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	return nil
}

func resourceColorCollectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	return nil
}