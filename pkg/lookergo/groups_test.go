package lookergo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroups_ListGroups(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/4.0/groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[
  {
    "can_add_to_content_metadata": true,
    "external_group_id": null,
    "id": "101",
    "name": "Admin jupiter-xcm",
    "user_count": 0,
    "externally_managed": false,
    "include_by_default": false,
    "contains_current_user": false,
    "can": {
      "show": true,
      "create": true,
      "index": true,
      "update": true,
      "delete": true,
      "edit_in_ui": true,
      "add_to_content_metadata": true
    }
  },
  {
    "can_add_to_content_metadata": true,
    "external_group_id": null,
    "id": "202",
    "name": "Admin saturn-wfs",
    "user_count": 0,
    "externally_managed": false,
    "include_by_default": false,
    "contains_current_user": false,
    "can": {
      "show": true,
      "create": true,
      "index": true,
      "update": true,
      "delete": true,
      "edit_in_ui": true,
      "add_to_content_metadata": true
    }
  }
]`)
	})

	groups, resp, err := client.Groups.List(ctx, nil)
	_ = resp
	if err != nil {
		t.Errorf("groups.List returned error: %v", err)
	}

	expectedGroups := []Group{
		{Id: 101, Name: "Admin jupiter-xcm", CanAddToContentMetadata: true},
		{Id: 202, Name: "Admin saturn-wfs", CanAddToContentMetadata: true}}
	if !reflect.DeepEqual(groups, expectedGroups) {
		t.Error(errGotWant("groups.List", groups, expectedGroups))
	}
}

func TestGroupsResourceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/4.0/groups/202", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
  "can_add_to_content_metadata": true,
  "external_group_id": null,
  "id": "202",
  "name": "Admin saturn-wfs",
  "user_count": 0,
  "externally_managed": false,
  "include_by_default": false,
  "contains_current_user": false,
  "can": {
    "show": true,
    "create": true,
    "index": true,
    "update": true,
    "delete": true,
    "edit_in_ui": true,
    "add_to_content_metadata": true
  }
}`)
	})

	group, resp, err := client.Groups.Get(ctx, 202)
	_ = resp
	if err != nil {
		t.Errorf("groups.Get returned error: %v", err)
	}

	expectedGroup := &Group{Id: 202, Name: "Admin saturn-wfs", CanAddToContentMetadata: true}
	if !reflect.DeepEqual(group, expectedGroup) {
		t.Error(errGotWant("groups.List", group, expectedGroup))
	}
}

func TestGroupsResourceOp_ListByName(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/4.0/groups/search/with_hierarchy", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("name") == "Demo GROUP" {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, `[{
			    "id": "101",
			    "name": "Demo GROUP",
			    "user_count": 1,
			    "parent_group_ids": ["99","44"],
			    "role_ids": ["33","34"]
			  }
			]`)
		} else if r.URL.Query().Get("name") == "Admin%" {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, `[{
			    "id": "1001",
			    "name": "Admin ONE",
			    "user_count": 1,
			    "parent_group_ids": ["88","99"],
			    "role_ids": ["44"]
			  },{
			    "id": "1002",
			    "name": "Admin TWO",
			    "user_count": 2,
			    "parent_group_ids": [],
			    "role_ids": ["55","99"]
			  },{
			    "id": "1003",
			    "name": "Admin THREE",
			    "user_count": 0,
			    "parent_group_ids": [],
			    "role_ids": []
			  }]`)
		} else {
			t.Errorf("Groups.ListByName did not request with a name parameter")
		}
	})

	t.Run("Search single entry", func(t *testing.T) {
		groups, resp, err := client.Groups.ListByName(ctx, "Demo GROUP", nil)
		_ = resp
		if err != nil {
			t.Errorf("groups.List returned error: %v", err)
		}

		expectedGroups := []Group{
			{Id: 101, Name: "Demo GROUP", UserCount: 1, ParentGroupIds: []int{99, 44}, RoleIds: []int{33, 34}}}
		if !reflect.DeepEqual(groups, expectedGroups) {
			t.Error(errGotWant("groups.ListByName", groups, expectedGroups))
		}
	})

	t.Run("Search wildcard multiple entries", func(t *testing.T) {
		groups, resp, err := client.Groups.ListByName(ctx, "Admin%", nil)
		_ = resp
		if err != nil {
			t.Errorf("groups.List returned error: %v", err)
		}

		expectedGroups := []Group{
			{Id: 1001, Name: "Admin ONE", UserCount: 1, ParentGroupIds: []int{88, 99}, RoleIds: []int{44}},
			{Id: 1002, Name: "Admin TWO", UserCount: 2, RoleIds: []int{55, 99}},
			{Id: 1003, Name: "Admin THREE", UserCount: 0}}
		if !reflect.DeepEqual(groups, expectedGroups) {
			t.Error(errGotWant("groups.ListByName", groups, expectedGroups))
		}
	})
}

func TestGroupsResourceOp_ListById(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/4.0/groups/search/with_hierarchy", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("id") == "101" {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, `[{
			    "id": "101",
			    "name": "Demo GROUP",
			    "user_count": 1,
			    "parent_group_ids": ["99","44"],
			    "role_ids": ["33","34"]
			  }
			]`)
		} else if r.URL.Query().Get("id") == "1001,1002,1003" {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, `[{
			    "id": "1001",
			    "name": "Admin ONE",
			    "user_count": 1,
			    "parent_group_ids": ["88","99"],
			    "role_ids": ["44"]
			  },{
			    "id": "1002",
			    "name": "Admin TWO",
			    "user_count": 2,
			    "parent_group_ids": [],
			    "role_ids": ["55","99"]
			  },{
			    "id": "1003",
			    "name": "Admin THREE",
			    "user_count": 0,
			    "parent_group_ids": [],
			    "role_ids": []
			  }]`)
		} else {
			t.Errorf("Groups.ListByName did not request with an id parameter")
		}
	})

	t.Run("Search single entry", func(t *testing.T) {
		groups, resp, err := client.Groups.ListById(ctx, []int{101}, nil)
		_ = resp
		if err != nil {
			t.Errorf("groups.List returned error: %v", err)
		}

		expectedGroups := []Group{
			{Id: 101, Name: "Demo GROUP", UserCount: 1, ParentGroupIds: []int{99, 44}, RoleIds: []int{33, 34}}}
		if !reflect.DeepEqual(groups, expectedGroups) {
			t.Error(errGotWant("groups.ListByName", groups, expectedGroups))
		}
	})

	t.Run("Search wildcard multiple entries", func(t *testing.T) {
		groups, resp, err := client.Groups.ListById(ctx, []int{1001, 1002, 1003}, nil)
		_ = resp
		if err != nil {
			t.Errorf("groups.List returned error: %v", err)
		}

		expectedGroups := []Group{
			{Id: 1001, Name: "Admin ONE", UserCount: 1, ParentGroupIds: []int{88, 99}, RoleIds: []int{44}},
			{Id: 1002, Name: "Admin TWO", UserCount: 2, RoleIds: []int{55, 99}},
			{Id: 1003, Name: "Admin THREE", UserCount: 0}}
		if !reflect.DeepEqual(groups, expectedGroups) {
			t.Error(errGotWant("groups.ListByName", groups, expectedGroups))
		}
	})
}

func TestGroupsResourceOp_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/4.0/groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{
  "can_add_to_content_metadata": true,
  "external_group_id": null,
  "id": "202",
  "name": "Admin saturn-wfs",
  "user_count": 0,
  "externally_managed": false,
  "include_by_default": false,
  "contains_current_user": false,
  "can": {
    "show": true,
    "create": true,
    "index": true,
    "update": true,
    "delete": true,
    "edit_in_ui": true,
    "add_to_content_metadata": true
  }
}`)
	})
	newEntry := &Group{Id: 202, Name: "Admin saturn-wfs", CanAddToContentMetadata: true}

	created, resp, err := client.Groups.Create(ctx, newEntry)
	_ = resp
	if err != nil {
		t.Errorf("groups.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(created, newEntry) {
		t.Error(errGotWant("groups.Create", created, newEntry))
	}

}

func TestGroupsResourceOp_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/4.0/groups/202", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		fmt.Fprint(w, `{
  "can_add_to_content_metadata": true,
  "external_group_id": null,
  "id": "202",
  "name": "Admin saturn-wfs",
  "user_count": 0,
  "externally_managed": false,
  "include_by_default": false,
  "contains_current_user": false,
  "can": {
    "show": true,
    "create": true,
    "index": true,
    "update": true,
    "delete": true,
    "edit_in_ui": true,
    "add_to_content_metadata": true
  }
}`)
	})
	updatedEntry := &Group{Id: 202, Name: "Admin saturn-wfs", CanAddToContentMetadata: true}

	created, resp, err := client.Groups.Update(ctx, 202, updatedEntry)
	_ = resp
	if err != nil {
		t.Errorf("groups.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(created, updatedEntry) {
		t.Error(errGotWant("groups.Update", created, updatedEntry))
	}
}

func TestGroupsResourceOp_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/4.0/groups/202", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Groups.Delete(ctx, 202)
	if err != nil {
		t.Errorf("groups.Delete returned error: %v", err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)
}

func TestGroups_ListMemberGroups(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/4.0/groups/99/groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[
  {
    "can_add_to_content_metadata": true,
    "external_group_id": null,
    "id": "101",
    "name": "Admin jupiter-xcm",
    "user_count": 0,
    "externally_managed": false,
    "include_by_default": false,
    "contains_current_user": false,
    "can": {
      "show": true,
      "create": true,
      "index": true,
      "update": true,
      "delete": true,
      "edit_in_ui": true,
      "add_to_content_metadata": true
    }
  },
  {
    "can_add_to_content_metadata": true,
    "external_group_id": null,
    "id": "202",
    "name": "Admin saturn-wfs",
    "user_count": 0,
    "externally_managed": false,
    "include_by_default": false,
    "contains_current_user": false,
    "can": {
      "show": true,
      "create": true,
      "index": true,
      "update": true,
      "delete": true,
      "edit_in_ui": true,
      "add_to_content_metadata": true
    }
  }
]`)
	})

	groups, resp, err := client.Groups.ListMemberGroups(ctx, 99, nil)
	_ = resp
	if err != nil {
		t.Errorf("groups.List returned error: %v", err)
	}

	expectedGroups := []Group{
		{Id: 101, Name: "Admin jupiter-xcm", CanAddToContentMetadata: true},
		{Id: 202, Name: "Admin saturn-wfs", CanAddToContentMetadata: true}}
	if !reflect.DeepEqual(groups, expectedGroups) {
		t.Error(errGotWant("groups.ListMemberGroups", groups, expectedGroups))
	}
}

func TestGroupsResourceOp_AddMemberGroup(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &NewGroupMemberGroup{GroupID: 202}
	createResponse := &Group{Id: 202, Name: "Admin saturn-wfs", CanAddToContentMetadata: true}

	mux.HandleFunc("/api/4.0/groups/999/groups", func(w http.ResponseWriter, r *http.Request) {
		v := new(NewGroupMemberGroup)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(v, createRequest) {
			t.Errorf("Request body = %+v, expected %+v", v, createRequest)
		}
		fmt.Fprint(w, `{
  "can_add_to_content_metadata": true,
  "external_group_id": null,
  "id": "202",
  "name": "Admin saturn-wfs",
  "user_count": 0,
  "externally_managed": false,
  "include_by_default": false,
  "contains_current_user": false,
  "can": {
    "show": true,
    "create": true,
    "index": true,
    "update": true,
    "delete": true,
    "edit_in_ui": true,
    "add_to_content_metadata": true
  }
}`)
	})

	item, resp, err := client.Groups.AddMemberGroup(ctx, 999, 202)
	_ = resp
	if err != nil {
		t.Errorf("groups.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(item, createResponse) {
		t.Error(errGotWant("groups.ListMemberGroups", item, createResponse))
	}
}

func TestGroupsResourceOp_RemoveMemberGroup(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/4.0/groups/202/groups/999", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Groups.RemoveMemberGroup(ctx, 202, 999)
	if err != nil {
		t.Errorf("groups.Delete returned error: %v", err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)
}
