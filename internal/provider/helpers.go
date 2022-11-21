package provider

import (
	"context"
	"fmt"
	"regexp"
	"runtime"
	"strconv"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func idAsString(val interface{}) string {
	var ret string
	switch v := val.(type) {
	case int:
		ret = strconv.Itoa(v)
	case string:
		if _, err := strconv.Atoi(v); err != nil {
			panic(err)
		}
		ret = v
	}
	return ret
}

func idAsInt(val interface{}) int {
	var ret int
	switch v := val.(type) {
	case int:
		ret = v
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		ret = id
	}

	return ret
}

func currFuncName() string {
	counter, _, _, success := runtime.Caller(1)

	if !success {
		return "functionName: runtime.Caller: failed"
	}
	re := regexp.MustCompile(`(?m)\.(\S*)$`)
	return re.FindStringSubmatch(runtime.FuncForPC(counter).Name())[1]
}

func castToPtr[T any](v T) *T {
	return &v
}

func boolPtr(b bool) *bool {
	return &b
}

func logTrace(ctx context.Context, msg string, additional ...any) {
	add := make(map[string]interface{})
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	file, line := details.FileLine(pc)
	if ok && details != nil {
		add["@caller"] = fmt.Sprintf("%s:%d", file, line)
		// add["@method"] = details.Name()
	}
	for i := 0; i < len(additional); i = i + 2 {
		key := additional[i]
		value := additional[i+1]
		// add[key.(string)] = value
		add[fmt.Sprintf("%s_gostring", key.(string))] = fmt.Sprintf("%#v", value)
	}
	for _, a := range additional {
		_ = a
	}
	tflog.Trace(ctx, msg, add)
}

func logErrDiag(ctx context.Context, diags diag.Diagnostics, msg string, varName string, varValue interface{}) diag.Diagnostics {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	file, line := details.FileLine(pc)
	if ok && details != nil {
		tflog.Error(ctx, msg, map[string]interface{}{
			varName:                             varValue,
			fmt.Sprintf("%s_gostring", varName): fmt.Sprintf("%#v", varValue),
			"@caller":                           fmt.Sprintf("%s:%d", file, line),
			"@method":                           details.Name(),
		})
		return append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       msg,
			Detail:        fmt.Sprintf("var %s (%[1]T): %#v\nfunc: %s\nin: %s:%d", varName, varValue, details.Name(), file, line),
			AttributePath: cty.Path{cty.GetAttrStep{Name: varName}},
		})
	}
	tflog.Error(ctx, msg, map[string]interface{}{
		varName: varValue,
	})
	return append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  msg,
		Detail:   fmt.Sprintf("var %s (%[1]T): %#v", varName, varValue),
	})
}

func getSetIds(d *schema.ResourceData, key string) (setIds []string) {
	itemSet := d.Get(key).(*schema.Set)
	for _, elemSet := range itemSet.List() {
		elemId := elemSet.(map[string]interface{})["id"].(string)
		setIds = append(setIds, elemId)
	}
	return
}

func getSetChangeIdsDiff(d *schema.ResourceData, key string) (old *schema.Set, new *schema.Set, oldIds []string, newIds []string) {
	oldRaw, newRaw := d.GetChange(key)

	old = oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set)) // .List()
	new = newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set)) // .List()

	for _, elem := range old.List() {
		elemId := elem.(map[string]interface{})["id"].(string)
		oldIds = append(oldIds, elemId)
	}
	for _, elem := range new.List() {
		elemId := elem.(map[string]interface{})["id"].(string)
		newIds = append(newIds, elemId)
	}

	return
}

func schemaSetToStringSlice(set1 *schema.Set) []string {
	ret := make([]string, set1.Len())
	for i, item := range set1.List() {
		ret[i] = item.(string)
	}
	return ret
}

/*func getMapIds(itemSet []interface{}) (ids []string) {
	for _, elem := range itemSet.List() {
		elemId := elemSet.(map[string]interface{})["id"].(string)
		setIds = append(setIds, elemId)
	}
	return
}
*/
