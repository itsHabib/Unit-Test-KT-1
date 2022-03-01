package cats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetBreeds(t *testing.T) {
	for _, tc := range []struct {
		desc    string
		wantErr bool
	}{
		{
			desc:    "get breeds requests failure should result in an error",
			wantErr: true,
		},
		{
			desc:    "get breeds requests should throw an error on a non-200 response",
			wantErr: true,
		},
		{
			desc: "happy path",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			c, err := NewClient("key")
			require.NoError(t, err)

			resp, err := c.GetBreeds("mcoo")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
