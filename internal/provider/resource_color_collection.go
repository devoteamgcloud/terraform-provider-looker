package provider

import (
	"context"
	//"fmt"
	//"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	//"github.com/hashicorp/terraform-plugin-log/tflog"
	"fmt"
	"regexp"

	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
							Description: "Unique identity string",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"label": {
							Description: "Label of palette",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"type": {
							Description: "Type of palette",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "Categorical",
						},
						"colors": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("#(?:[0-9a-f]{2}){2,4}")
										return ret
									}(),
									"color must be a valid color hex code",
								)),
							},
						},
					},
				},
			},
			"sequentialPalettes": {
				Description: "Array of categorical palette definitions",
				Type:        schema.TypeList,
				Optional:    true,
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
							Default:     "Sequential",
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
										ValidateDiagFunc: validation.ToDiagFunc(
											validation.StringMatch(
												func() *regexp.Regexp {
													ret, _ := regexp.Compile("#(?:[0-9a-f]{2}){2,4}")
													return ret
												}(),
												"color must be a valid color hex code",
											),
										),
									},
									"offset": {
										Type:             schema.TypeInt,
										Description:      "Offset in continuous palette (0 to 100)",
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
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
				Optional:    true,
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
							Default:     "Diverging",
						},
						"stops": {
							Type:        schema.TypeList,
							Description: "Array of ColorStops in the palette",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"color": {
										Type:        schema.TypeString,
										Description: "CSS color string",
										Required:    true,
										ValidateFunc: validation.StringMatch(
											func() *regexp.Regexp {
												ret, _ := regexp.Compile("#(?:[0-9a-f]{2}){2,4}")
												return ret
											}(),
											"color must be a valid color hex code",
										),
									},
									"offset": {
										Type:             schema.TypeInt,
										Description:      "Offset in continuous palette (0 to 100)",
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
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
	// Checks whether the API Client is configured. If not, the resource responds with an error.
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	// Retrieves values from the plan. The function will attempt to retrieve values from the plan and convert it to an WriteColorCollection
	coco := &lookergo.WriteColorCollection{}

	coco.Label = castToPtr(d.Get("label").(string)) // schema.TypeString -> string

	schemaToDiscretePalette := func(str string) *[]lookergo.DiscretePalette {
		if catSchema, ok := d.GetOk(str); ok {
			catPal := []lookergo.DiscretePalette{}

			// schema.TypeList -> []interface{}
			for _, elem := range catSchema.([]interface{}) {
				e := elem.(map[string]interface{}) // map[string]*schema.Schema -> map[string]interface{}
				var cat lookergo.DiscretePalette   // TODO is this better out of the loop ?
				cat.Label = castToPtr(e["label"].(string))
				// https://stackoverflow.com/questions/59714262/type-assertion-for-typelist-in-terraform-provider
				rawColors := e["colors"].([]interface{})
				items := make([]string, len(rawColors))
				for i, raw := range rawColors {
					items[i] = raw.(string)
				}
				cat.Colors = &items

				catPal = append(catPal, cat)
			}

			return &catPal
		} else {
			return nil
		}
	}

	coco.CategoricalPalettes = schemaToDiscretePalette("categoricalPalettes")

	schemaToContinousPalette := func(str string) *[]lookergo.ContinuousPalette {
		if catSchema, ok := d.GetOk(str); ok {
			catPal := []lookergo.ContinuousPalette{}

			// schema.TypeList -> []interface{}
			for _, elem := range catSchema.([]interface{}) {
				e := elem.(map[string]interface{}) // map[string]*schema.Schema -> (map[string]interface{}
				var cat lookergo.ContinuousPalette // TODO is this better out of the loop ?
				cat.Label = castToPtr(e["label"].(string))
				// https://stackoverflow.com/questions/59714262/type-assertion-for-typelist-in-terraform-provider
				rawStops := e["stops"].([]interface{})
				items := make([]lookergo.ColorStop, len(rawStops))
				for i, raw := range rawStops {
					rawStop := raw.(map[string]interface{})
					items[i] = lookergo.ColorStop{
						Color:  castToPtr(rawStop["color"].(string)),
						Offset: castToPtr(rawStop["offset"].(int64)),
					}
				}
				cat.Stops = &items

				catPal = append(catPal, cat)
			}

			return &catPal
		} else {
			return nil
		}
	}

	coco.SequentialPalettes = schemaToContinousPalette("sequentialPalettes")
	coco.DivergingPalettes = schemaToContinousPalette("divergingPalettes")

	// send POST request. Creates a new order. The function invokes the API client's create method.
	newCoCo, _, err := c.ColorCollection.Create(ctx, coco)
	if err != nil {
		return diag.FromErr(err)
	}

	// Map response body to schema and populate Computed attribute values
	// TODO set IDs of palettes
	d.SetId(*newCoCo.Id)

	/*
		// regarding nested computed fields:
		// https://github.com/hashicorp/terraform/issues/10532
		if catId, ok := d.GetOk("categoricalPalettes"); ok {
			catPalRes := catId.(schema.ResourceData)
			catPalRes.Set()
			for _, c := range catId.([]interface{}) {
				c.set
			}
		}
	*/

	// check if resource has been created correctly
	// populate fields in ColorCollection with the values from resourceColorCollection
	return resourceColorCollectionRead(ctx, d, m)
}

func resourceColorCollectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	cocoID := d.Id()

	coco, _, err := c.ColorCollection.Get(ctx, cocoID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("id", *coco.Id); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("label", coco.Label); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("categoricalPalettes", *coco.CategoricalPalettes); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("sequentialPalettes", *coco.SequentialPalettes); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("divergingPalettes", *coco.DivergingPalettes); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceColorCollectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	cocoID := d.Id()

	coco, _, err := c.ColorCollection.Get(ctx, cocoID)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges() {
		if err = d.Set("id", *coco.Id); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("label", coco.Label); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("categoricalPalettes", *coco.CategoricalPalettes); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("sequentialPalettes", *coco.SequentialPalettes); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("divergingPalettes", *coco.DivergingPalettes); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceColorCollectionRead(ctx, d, m)
}

func resourceColorCollectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	CoCoId := d.Id()

	if _, err := c.ColorCollection.Delete(ctx, CoCoId); err != nil {
		return diag.FromErr(err)
	}
	// Finally mark as deleted
	d.SetId("")

	return diags
}
