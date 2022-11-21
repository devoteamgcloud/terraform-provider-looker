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
				Required:    true,
			},
			"categoricalpalettes": {
				Description: "Array of categorical palette definitions",
				Type:        schema.TypeSet,
				Required:    true,
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
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"sequentialpalettes": {
				Description: "Array of categorical palette definitions",
				Type:        schema.TypeSet,
				Required:    true,
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
							Type:        schema.TypeSet,
							Required:    true,
							MinItems:    2,
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
													ret, _ := regexp.Compile("#(?i)(?:[0-9A-F]{2}){2,4}")
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
			"divergingpalettes": {
				Description: "Array of categorical palette definitions",
				Type:        schema.TypeSet,
				Required:    true,
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
							Type:        schema.TypeSet,
							Description: "Array of ColorStops in the palette",
							Required:    true,
							MinItems:    2,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"color": {
										Type:        schema.TypeString,
										Description: "CSS color string",
										Required:    true,
										ValidateFunc: validation.StringMatch(
											func() *regexp.Regexp {
												ret, _ := regexp.Compile("#(?i)(?:[0-9A-F]{2}){2,4}")
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

func cocoSchemaToStruct(ctx context.Context, d *schema.ResourceData, coco *lookergo.WriteColorCollection) {
	if cocoLabel, ok := d.GetOk("label"); ok {
		coco.Label = castToPtr(cocoLabel.(string))
	}

	if categoricalpalettesSet, ok := d.GetOk("categoricalpalettes"); ok {
		catPal := []lookergo.DiscretePalette{}
		for _, raw := range categoricalpalettesSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			var pal lookergo.DiscretePalette
			pal.Label = castToPtr(obj["label"].(string))
			pal.Colors = castToPtr(schemaSetToStringSlice(obj["colors"].(*schema.Set)))
			pal.Type = castToPtr(obj["type"].(string))
			catPal = append(catPal, pal)
		}
		coco.CategoricalPalettes = &catPal
	}

	if sequentialpalettesSet, ok := d.GetOk("sequentialpalettes"); ok {
		seqPal := []lookergo.ContinuousPalette{}
		for _, raw := range sequentialpalettesSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			var pal lookergo.ContinuousPalette
			pal.Label = castToPtr(obj["label"].(string))
			pal.Type = castToPtr(obj["type"].(string))
			stopsList := []lookergo.ColorStop{}
			for _, stop := range obj["stops"].(*schema.Set).List() {
				objj := stop.(map[string]interface{})
				st := lookergo.ColorStop{}
				st.Color = castToPtr(objj["color"].(string))
				st.Offset = castToPtr(objj["offset"].(int))
				stopsList = append(stopsList, st)
			}
			pal.Stops = &stopsList
			seqPal = append(seqPal, pal)
		}
		coco.SequentialPalettes = &seqPal
	}

	if divergingpalettesSet, ok := d.GetOk("divergingpalettes"); ok {
		divPal := []lookergo.ContinuousPalette{}
		for _, raw := range divergingpalettesSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			var pal lookergo.ContinuousPalette
			pal.Label = castToPtr(obj["label"].(string))
			pal.Type = castToPtr(obj["type"].(string))
			stopsList := []lookergo.ColorStop{}
			for _, stop := range obj["stops"].(*schema.Set).List() {
				objj := stop.(map[string]interface{})
				st := lookergo.ColorStop{}
				st.Color = castToPtr(objj["color"].(string))
				st.Offset = castToPtr(objj["offset"].(int))
				//return diag.Errorf(objj["offset"].(string))
				stopsList = append(stopsList, st)
			}
			pal.Stops = &stopsList
			divPal = append(divPal, pal)
		}
		coco.DivergingPalettes = &divPal
	}
}

// Receives terraform resource schema, builds a golang struct with json fields from it, sends a Post request with the
func resourceColorCollectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	// Checks whether the API Client is configured. If not, the resource responds with an error.
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	// Retrieves values from the plan. The function will attempt to retrieve values from the plan and convert it to an WriteColorCollection
	var coco lookergo.WriteColorCollection

	cocoSchemaToStruct(ctx, d, &coco)

	// send POST request. Creates a new order. The function invokes the API client's create method.
	newCoCo, _, err := c.ColorCollection.Create(ctx, &coco)
	if err != nil {
		fmt.Println(newCoCo)
		return diag.FromErr(err)
	}
	d.SetId(*newCoCo.Id)
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
	if err = d.Set("label", *coco.Label); err != nil {
		return diag.FromErr(err)
	}

	flattenDiscretePalette := func(coco *[]lookergo.DiscretePalette) []interface{} {
		catPals := make([]interface{}, len(*coco), len(*coco))
		for i, elem := range *coco {
			catPal := make(map[string]interface{})
			catPal["label"] = *elem.Label
			catPal["id"] = *elem.Id
			catPal["type"] = *elem.Type
			catPal["colors"] = *elem.Colors
			catPals[i] = catPal
		}
		return catPals
	}
	catPals := flattenDiscretePalette(coco.CategoricalPalettes)
	if err = d.Set("categoricalpalettes", catPals); err != nil {
		return diag.FromErr(err)
	}
	flattenContinuousPalette := func(coco *[]lookergo.ContinuousPalette) []interface{} {
		contPals := make([]interface{}, len(*coco), len(*coco))
		for i, elem := range *coco {
			xPal := make(map[string]interface{})
			xPal["label"] = *elem.Label
			xPal["id"] = *elem.Id
			xPal["type"] = *elem.Type
			// Flatten stops
			stops := make([]map[string]interface{}, len(*elem.Stops), len(*elem.Stops))
			for i, stop := range *elem.Stops {
				s := make(map[string]interface{})
				s["color"] = *stop.Color
				s["offset"] = *stop.Offset
				stops[i] = s
			}
			xPal["stops"] = stops
			contPals[i] = xPal
		}

		return contPals

	}

	seqPals := flattenContinuousPalette(coco.SequentialPalettes)
	if err = d.Set("sequentialpalettes", seqPals); err != nil {
		return diag.FromErr(err)
	}

	divPals := flattenContinuousPalette(coco.DivergingPalettes)
	if err = d.Set("divergingpalettes", divPals); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceColorCollectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	cocoID := d.Id()

	if d.HasChanges("label", "catogoricalpalletes", "sequentialpalettes", "divergingpalettes") {
		var coco lookergo.WriteColorCollection
		cocoSchemaToStruct(ctx, d, &coco)
		newCoco, _, err := c.ColorCollection.Update(ctx, cocoID, &coco)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("id", *newCoco.Id)
		d.SetId(*newCoco.Id)
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
