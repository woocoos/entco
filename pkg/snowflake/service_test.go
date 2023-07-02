package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsingsun/woocoo/pkg/conf"
	"os"
	"testing"
)

func TestSetDefaultNode(t *testing.T) {
	type args struct {
		cnf *conf.Configuration
	}
	tests := []struct {
		name  string
		args  args
		panic bool
		check func()
	}{
		{
			name: "init",
			args: args{
				cnf: nil,
			},
			check: func() {
				assert.Equal(t, uint8(10), snowflake.NodeBits)
				assert.Equal(t, uint8(12), snowflake.StepBits)
				id := New()
				assert.Len(t, id.String(), 19)
			},
		},
		{
			name: "default",
			args: args{
				cnf: conf.NewFromStringMap(map[string]any{}),
			},
			panic: false,
			check: func() {
				id := New()
				assert.Len(t, id.String(), 19)
			},
		},
		{
			name: "small",
			args: args{
				cnf: func() *conf.Configuration {
					return conf.NewFromStringMap(map[string]any{
						"nodeBits": 1,
						"stepBits": 8,
					})
				}(),
			},
			panic: false,
			check: func() {
				assert.Equal(t, uint8(1), snowflake.NodeBits)
				assert.Equal(t, uint8(8), snowflake.StepBits)
				id := New()
				assert.Len(t, id.String(), 15)
			},
		},
		{
			name: "node from env",
			args: args{
				cnf: func() *conf.Configuration {
					defer func() {
						os.Setenv("SNOWFLAKE_NODE_ID", "")
					}()
					require.NoError(t, os.Setenv("SNOWFLAKE_NODE_ID", "2"))
					return conf.NewFromStringMap(map[string]any{
						"nodeBits": 2,
						"stepBits": 8,
					})
				}(),
			},
			panic: false,
			check: func() {
				assert.Equal(t, uint8(2), snowflake.NodeBits)
				assert.Equal(t, uint8(8), snowflake.StepBits)
				id := New()
				assert.Len(t, id.String(), 15)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				assert.Error(t, SetDefaultNode(tt.args.cnf))
				return
			}
			if tt.args.cnf != nil {
				require.NoError(t, SetDefaultNode(tt.args.cnf))
			}
			tt.check()
		})
	}
}
