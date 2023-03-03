package initx

import (
	"entgo.io/ent/dialect"
	"github.com/tsingsun/woocoo/pkg/conf"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestBuildEntCacheDriver(t *testing.T) {
	type args struct {
		cnf       *conf.AppConfiguration
		preDriver dialect.Driver
	}
	tests := []struct {
		name string
		args args
		want dialect.Driver
	}{
		{
			name: "no set",
			args: args{
				cnf: &conf.AppConfiguration{
					Configuration: conf.New(),
				},
				preDriver: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildEntCacheDriver(tt.args.cnf, tt.args.preDriver); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildEntCacheDriver() = %v, want %v", got, tt.want)
			}
		})
	}
}
