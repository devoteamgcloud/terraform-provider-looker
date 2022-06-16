package provider

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
)

func idAsString(val interface{}) string {
	var ret string
	switch val.(type) {
	case int:
		ret = strconv.Itoa(val.(int))
	case string:
		if _, err := strconv.Atoi(val.(string)); err != nil {
			panic(err)
		}
		ret = val.(string)
	}

	return ret
}

func idAsInt(val interface{}) int {
	var ret int
	switch val.(type) {
	case int:
		ret = val.(int)
	case string:
		id, err := strconv.Atoi(val.(string))
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
		return fmt.Sprintf("functionName: runtime.Caller: failed")
	}
	re := regexp.MustCompile(`(?m)\.(\S*)$`)
	return re.FindStringSubmatch(runtime.FuncForPC(counter).Name())[1]
}

func boolPtr(b bool) *bool {
	return &b
}

/*
func multiSchemaSet(d *schema.ResourceData, diags diag.Diagnostics, x map[string]interface{}) {
	// field :=
	for key, element := range x {
		if err := d.Set(key, element); err != nil {
			diag.FromErr(err)

		}
	}
}
*/
//  end
