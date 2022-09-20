package dmap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_lookup(t *testing.T) {
	type args struct {
		prefix     string
		sep        string
		sourceData map[string]interface{}
		destData   map[string]interface{}
	}
	type want struct {
		want map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "ok",
			args: args{
				prefix: "mysql",
				sep:    "_",
				sourceData: map[string]interface{}{
					"dev": map[string]interface{}{
						"host": "127.0.0.1",
						"port": 3306,
					},
					"stag": map[string]interface{}{
						"host": "https://stag.qa.com",
						"port": 3306,
					},
				},
				destData: make(map[string]interface{}, 2),
			},
			want: want{want: map[string]interface{}{
				"mysql_dev_host":  "127.0.0.1",
				"mysql_dev_port":  3306,
				"mysql_stag_host": "https://stag.qa.com",
				"mysql_stag_port": 3306,
			}},
		},
		{
			name: "without prefix",
			args: args{
				prefix: "",
				sep:    ".",
				sourceData: map[string]interface{}{
					"dev": map[string]interface{}{
						"host": "127.0.0.1",
						"port": 3306,
					},
					"stag": map[string]interface{}{
						"host": "https://stag.qa.com",
						"port": 3306,
					},
				},
				destData: make(map[string]interface{}, 2),
			},
			want: want{want: map[string]interface{}{
				"dev.host":  "127.0.0.1",
				"dev.port":  3306,
				"stag.host": "https://stag.qa.com",
				"stag.port": 3306,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lookup(tt.args.prefix, tt.args.sep, tt.args.sourceData, tt.args.destData)
			assert.Equal(t, tt.want.want, tt.args.destData)
		})
	}
}
