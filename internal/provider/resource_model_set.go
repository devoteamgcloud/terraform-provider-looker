package provider

import (
	"context"
	"fmt"

	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceModelSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModelSetCreate,
		ReadContext:   resourceModelSetRead,
		UpdateContext: resourceModelSetUpdate,
		DeleteContext: resourceModelSetDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "Name for ModelSet of LookML Models",
				Type:        schema.TypeString,
				Required:    true,
			},
			"models": {
				Description: "List of LookML Model names",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceModelSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))
	var interfaceModel = d.Get("models").(*schema.Set)
	var model_list []string
	for _, raw := range interfaceModel.List() {
		model_list = append(model_list, raw.(string))
	}
	var modelSet = lookergo.ModelSet{Name: d.Get("name").(string), Models: model_list}
	newSet, _, err := c.ModelSets.Create(ctx, &modelSet)
	if err != nil {
		return diag.FromErr(err)
	}
	if newSet == nil {
		return diag.FromErr(&lookergo.ArgError{})
	}
	d.SetId(newSet.Id)
	resourceModelSetRead(ctx, d, m)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceModelSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))
	var id = d.Id()
	newModel, _, err := c.ModelSets.Get(ctx, id)

	if err != nil {
		return diag.FromErr(err)
	}
	if newModel == nil {
		return diag.FromErr(&lookergo.ArgError{})
	}
	if err = d.Set("id", newModel.Id); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", newModel.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("models", newModel.Models); err != nil {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceModelSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	// TODO
	_ = c

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return resourceModelSetRead(ctx, d, m)
}

func resourceModelSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	// TODO
	_ = c

	// Finally mark as deleted
	d.SetId("")
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}
