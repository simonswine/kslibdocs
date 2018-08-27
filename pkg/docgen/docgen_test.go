package docgen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_sortVersions(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		out  []string
	}{
		{
			name: "one version",
			in:   []string{"v1"},
			out:  []string{"v1"},
		},
		{
			name: "multiple versions 1",
			in:   []string{"v1beta1", "v1beta2", "v1"},
			out:  []string{"v1", "v1beta2", "v1beta1"},
		},
		{
			name: "multiple versions 2",
			in:   []string{"v1beta1", "v1"},
			out:  []string{"v1", "v1beta1"},
		},
		{
			name: "multiple versions 3",
			in:   []string{"v1", "v1beta1"},
			out:  []string{"v1", "v1beta1"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := sortVersions(tc.in)
			require.Equal(t, tc.out, out)
		})
	}
}
