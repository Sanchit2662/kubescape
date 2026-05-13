package opaprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCELEvaluator(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		resource    map[string]interface{}
		wantMatch   bool
		wantErr     bool
	}{
		{
			name:       "privileged container detected",
			expression: `object.spec.containers.exists(c, c.securityContext.privileged == true)`,
			resource: map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name": "app",
							"securityContext": map[string]interface{}{
								"privileged": true,
							},
						},
					},
				},
			},
			wantMatch: true,
		},
		{
			name:       "non-privileged container passes",
			expression: `object.spec.containers.exists(c, c.securityContext.privileged == true)`,
			resource: map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name": "app",
							"securityContext": map[string]interface{}{
								"privileged": false,
							},
						},
					},
				},
			},
			wantMatch: false,
		},
		{
			name:       "malformed expression returns error",
			expression: `this is not valid CEL !!!`,
			resource:   map[string]interface{}{},
			wantErr:    true,
		},
		{
			name:       "non-bool expression returns error",
			expression: `"hello"`,
			resource:   map[string]interface{}{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CELEvaluator(tt.expression, tt.resource)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantMatch, got)
		})
	}
}
