package docgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_kindSummary(t *testing.T) {
	group := "group"
	kind := "kind"
	versions := []string{"v2", "v1"}

	descriptions := map[string]string{
		"v2": "This is v2.",
		"v1": "This is v1.",
	}

	summary, err := kindSummary(group, kind, versions, descriptions)
	require.NoError(t, err)

	expected := "\n<div id=\"group-v2-kind\" class=\"kind-summary default\">\n\tThis is v2.\n</div>\n\n<div id=\"group-v1-kind\" class=\"kind-summary\">\n\tThis is v1.\n</div>\n\n"
	assert.Equal(t, expected, summary)
}
