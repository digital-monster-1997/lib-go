package dmap

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"reflect"
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
func Test_DeepSearchInMap(t *testing.T) {
	type args struct {
		sourceMap map[string]interface{}
		paths     []string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "ok",
			args: args{
				paths: []string{"stag"}, // Flattened multilevel map deep search，
				sourceMap: map[string]interface{}{
					"stag": map[string]interface{}{
						"host": "https://stag.qa.com",
						"port": 3306,
					},
				},
			},
			want: map[string]interface{}{
				"host": "https://stag.qa.com",
				"port": 3306,
			},
		},
		{
			name: "failed to get data",
			args: args{
				paths: []string{"dev"}, // Flattened multilevel map deep search，
				sourceMap: map[string]interface{}{
					"stag": map[string]interface{}{
						"host": "https://stag.qa.com",
						"port": 3306,
					},
				},
			},
			want: map[string]interface{}{},
		},
		{
			name: "go too deeper to  get data",
			args: args{
				paths: []string{"stag", "host"}, // Flattened multilevel map deep search，
				sourceMap: map[string]interface{}{
					"stag": map[string]interface{}{
						"host": "https://stag.qa.com",
						"port": 3306,
					},
				},
			},
			want: map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, DeepSearchInMap(tt.args.sourceMap, tt.args.paths...), "DeepSearchInMap(%v, %v)", tt.args.sourceMap, tt.args)
		})
	}
}

func TestToMapStringInterface(t *testing.T) {
	type args struct {
		src map[interface{}]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "convert int to string",
			args: args{src: map[interface{}]interface{}{
				1: "value",
			}},
			want: map[string]interface{}{
				"1": "value",
			},
		},
		{
			name: "convert bool to string",
			args: args{src: map[interface{}]interface{}{
				true: "value",
			}},
			want: map[string]interface{}{
				"true": "value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ToMapStringInterface(tt.args.src), "ToMapStringInterface(%v)", tt.args.src)
		})
	}
}

func TestMergeStringMap(t *testing.T) {
	type args struct {
		dest map[string]interface{}
		src  map[string]interface{}
		tar  map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "二維測試",
			args: args{
				dest: map[string]interface{}{
					"2w": map[string]interface{}{
						"test":  "2wtd",
						"test1": "2wtd1",
					},
					"2wa": map[string]interface{}{
						"test":  "2wtd",
						"test1": "2wtd1",
					},
					"2wi": map[interface{}]interface{}{
						"test":  "2wtd",
						"test1": "2wtd1",
					},
				},
				src: map[string]interface{}{
					"2w": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
					"2wb": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
					"2wi": map[interface{}]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
				},
				tar: map[string]interface{}{
					"2w": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
					"2wb": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
					"2wa": map[string]interface{}{
						"test":  "2wtd",
						"test1": "2wtd1",
					},
					"2wi": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
				},
			},
		},
		{
			name: "ok",
			args: args{
				dest: map[string]interface{}{
					"1w":  "tt",
					"1wa": "mq",
				},
				src: map[string]interface{}{
					"1w":  "tts",
					"1wb": "bq",
				},
				tar: map[string]interface{}{
					"1w":  "tts",
					"1wa": "mq",
					"1wb": "bq",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MergeStringMap(tt.args.dest, tt.args.src)
			if !reflect.DeepEqual(tt.args.dest, tt.args.tar) {
				spew.Dump(tt.args.dest)
				t.FailNow()
			}
		})
	}
}
