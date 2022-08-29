package lookergo

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRolesResourceOp_ModelSetsList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/4.0/model_sets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[
  { "built_in": true, "id": "1", "all_access": true, "models": ["another_dink", "dink_th_test"], "name": "All", "url": "https://localhost:19999/api/4.0/model_sets/1", "can": { "show": true, "index": true, "update": true } },
  { "built_in": false, "id": "2", "all_access": false, "models": ["accounts", "documents"], "name": "Model Set", "url": "https://localhost:19999/api/4.0/model_sets/2", "can": { "show": true, "index": true, "update": true } },
  { "built_in": false, "id": "5", "all_access": false, "models": ["dink_th_test"], "name": "Dink TH", "url": "https://localhost:19999/api/4.0/model_sets/5", "can": { "show": true, "index": true, "update": true } }
]`)
	})

	result, resp, err := client.ModelSets.List(ctx)
	_ = resp
	if err != nil {
		t.Errorf("Projects.Get returned error: %v", err)
	}

	expected := []ModelSet{
		{
			BuiltIn:   true,
			Id:        "1",
			AllAccess: true,
			Models:    []string{"another_dink", "dink_th_test"},
			Name:      "All",
		},
		{Id: "2", Models: []string{"accounts", "documents"}, Name: "Model Set"},
		{Id: "5", Models: []string{"dink_th_test"}, Name: "Dink TH"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Error(errGotWant("ModelSets.List", result, expected))
	}

}
