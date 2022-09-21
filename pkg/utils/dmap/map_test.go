package dmap

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestNewFlatMap(t *testing.T) {
	tests := []struct {
		name string
		want *FlatMap
	}{
		{
			name: "ok",
			want: &FlatMap{
				data: map[string]interface{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewFlatMap())
		})
	}
}

func TestFlatMap_Set(t *testing.T) {
	type args struct {
		key string
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "ok",
			args: args{key: "set_first_value", val: true},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewFlatMap()
			err := m.Set(tt.args.key, tt.args.val)
			assert.NoError(t, err)
			result := m.Get(tt.args.key)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestFlatMap_GetString(t *testing.T) {
	type args struct {
		key    string
		value  string
		getKey string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{key: "goodKey", value: "Good To See You", getKey: "goodKey"},
			want: "Good To See You",
		},
		{
			name: "failed to get value",
			args: args{key: "goodKey", value: "Good To See You", getKey: "badKey"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewFlatMap()
			err := m.Set(tt.args.key, tt.args.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, m.GetString(tt.args.getKey))
		})
	}
}

func TestFlatMap_Load(t *testing.T) {
	type args struct {
		content   []byte
		unmarshal Unmarshaler
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			args: args{
				content: []byte(`
mysql:
  host: https://good.db.com
  port: 3306

`),
				unmarshal: yaml.Unmarshal,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return false
			},
		},
		{
			name: "failed to unmarshal yaml",
			args: args{
				content:   []byte("!!!===!!!"),
				unmarshal: yaml.Unmarshal,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, err.Error(), "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!!===!!! `` into map[string]interface {}")
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flat := NewFlatMap()
			tt.wantErr(t, flat.Load(tt.args.content, tt.args.unmarshal))
		})
	}
}
