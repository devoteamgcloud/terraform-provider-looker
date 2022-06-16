package acceptance_tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAcptGroups_ListGroups(t *testing.T) {
	setup()

	groups, _, err := client.Groups.List(ctx, nil)
	if err != nil {
		t.Errorf("Groups.List returned error: %v", err)
	}

	assert.NotZero(t, groups)
}

func TestAcptGroups_GetGroup(t *testing.T) {
	setup()

	groups, _, err := client.Groups.Get(ctx, 1)
	if err != nil {
		t.Errorf("Groups.List returned error: %v", err)
	}

	assert.NotZero(t, groups)
}
