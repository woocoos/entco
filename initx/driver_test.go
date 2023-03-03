package initx

import (
	"entgo.io/ent/dialect"
	"github.com/stretchr/testify/assert"
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

func TestBuildOTELDriver(t *testing.T) {
	type args struct {
		cnf      *conf.AppConfiguration
		storekey string
	}
	tests := []struct {
		name  string
		args  args
		check func(driver dialect.Driver)
	}{
		{
			name: "otel",
			args: args{
				cnf: &conf.AppConfiguration{
					Configuration: conf.NewFromStringMap(map[string]any{
						"otel": "",
						"store": map[string]any{
							"mysql": map[string]any{
								"driverName": "mysql",
								"dsn":        "root:@tcp(localhost:3306)/portal?parseTime=true&loc=Local",
							},
						},
					}),
				},
				storekey: "store.mysql",
			},
			check: func(driver dialect.Driver) {
				assert.NotNil(t, driver)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildOTELDriver(tt.args.cnf, tt.args.storekey)
			tt.check(got)
		})
	}
}
