package generator

import (
	"context"
	"testing"

	"github.com/neonmei/szgen/internal/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		value   string
		count   int
		wantErr bool
	}{
		{
			name:    "constant",
			pattern: consts.GeneratorConstant,
			value:   "1",
			count:   1,
			wantErr: false,
		},
		{
			name:    "random",
			pattern: consts.GeneratorRandom,
			value:   "1,10", // min, max
			count:   1,
			wantErr: false,
		},
		{
			name:    "step",
			pattern: consts.GeneratorStep,
			value:   "0",
			count:   1,
			wantErr: false,
		},
		{
			name:    "sine",
			pattern: consts.GeneratorSine,
			value:   "10,10,0",
			count:   1,
			wantErr: false,
		},
		{
			name:    "sequence",
			pattern: consts.GeneratorSequence,
			value:   "1,2,3",
			count:   1,
			wantErr: false,
		},
		{
			name:    "unknown pattern",
			pattern: "unknown",
			value:   "",
			count:   1,
			wantErr: true,
		},
		{
			name:    "error in specific generator",
			pattern: consts.GeneratorConstant,
			value:   "invalid", // Invalid int64
			count:   1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := New[int64](context.Background(), tt.pattern, tt.value, tt.count)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, gen)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, gen)
			}
		})
	}
}
