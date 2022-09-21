package dmap

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Unmarshaler ...
type Unmarshaler = func([]byte, interface{}) error

// KeySplit ...
var KeySplit = "."

// FlatMap ...

type FlatMap struct {
	data   map[string]interface{}
	mu     sync.RWMutex
	keyMap sync.Map
}

// NewFlatMap ...
func NewFlatMap() *FlatMap {
	return &FlatMap{
		data: make(map[string]interface{}),
	}
}

// Load ...
func (flat *FlatMap) Load(content []byte, unmarshal Unmarshaler) error {
	data := make(map[string]interface{})
	if err := unmarshal(content, &data); err != nil {
		return err
	}
	return flat.apply(data)
}

func (flat *FlatMap) apply(data map[string]interface{}) error {
	flat.mu.Lock()
	defer flat.mu.Unlock()

	MergeStringMap(flat.data, data)
	var changes = make(map[string]interface{})
	for k, v := range flat.traverse(KeySplit) {
		orig, ok := flat.keyMap.Load(k)
		if ok && !reflect.DeepEqual(orig, v) {
			changes[k] = v
		}
		flat.keyMap.Store(k, v)
	}

	return nil
}

func (flat *FlatMap) traverse(sep string) map[string]interface{} {
	data := make(map[string]interface{})
	lookup("", sep, data, flat.data)
	return data
}

// Set ...
func (flat *FlatMap) Set(key string, val interface{}) error {
	paths := strings.Split(key, KeySplit)
	lastKey := paths[len(paths)-1]
	m := deepSearch(flat.data, paths[:len(paths)-1])
	m[lastKey] = val
	return flat.apply(m)
}

// Get returns the value associated with the key
func (flat *FlatMap) Get(key string) interface{} {
	return flat.find(key)
}

// GetString returns the value associated with the key as a string.
func (flat *FlatMap) GetString(key string) string {
	return cast.ToString(flat.Get(key))
}

// GetBool returns the value associated with the key as a boolean.
func (flat *FlatMap) GetBool(key string) bool {
	return cast.ToBool(flat.Get(key))
}

// GetInt returns the value associated with the key as an integer.
func (flat *FlatMap) GetInt(key string) int {
	return cast.ToInt(flat.Get(key))
}

// GetInt64 returns the value associated with the key as an integer.
func (flat *FlatMap) GetInt64(key string) int64 {
	return cast.ToInt64(flat.Get(key))
}

// GetFloat64 returns the value associated with the key as a float64.
func (flat *FlatMap) GetFloat64(key string) float64 {
	return cast.ToFloat64(flat.Get(key))
}

// GetTime returns the value associated with the key as time.
func (flat *FlatMap) GetTime(key string) time.Time {
	return cast.ToTime(flat.Get(key))
}

// GetDuration returns the value associated with the key as a duration.
func (flat *FlatMap) GetDuration(key string) time.Duration {
	return cast.ToDuration(flat.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (flat *FlatMap) GetStringSlice(key string) []string {
	return cast.ToStringSlice(flat.Get(key))
}

// GetSlice returns the value associated with the key as a slice of strings.
func (flat *FlatMap) GetSlice(key string) []interface{} {
	return cast.ToSlice(flat.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (flat *FlatMap) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(flat.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (flat *FlatMap) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(flat.Get(key))
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (flat *FlatMap) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(flat.Get(key))
}

// UnmarshalKey takes a single key and unmarshal it into a Struct.
func (flat *FlatMap) UnmarshalKey(key string, rawVal interface{}, tagName string) error {
	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     rawVal,
		TagName:    tagName,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	if key == "" {
		flat.mu.RLock()
		defer flat.mu.RUnlock()
		return decoder.Decode(flat.data)
	}

	value := flat.Get(key)
	if value == nil {
		return fmt.Errorf("invalid key %s, maybe not exist in config", key)
	}

	return decoder.Decode(value)
}

// Reset ...
func (flat *FlatMap) Reset() {
	flat.mu.Lock()
	defer flat.mu.Unlock()

	flat.data = make(map[string]interface{})
	// erase map
	flat.keyMap.Range(func(key interface{}, value interface{}) bool {
		flat.keyMap.Delete(key)
		return true
	})
}

func (flat *FlatMap) find(key string) interface{} {
	dd, ok := flat.keyMap.Load(key)
	if ok {
		return dd
	}

	paths := strings.Split(key, KeySplit)
	flat.mu.RLock()
	defer flat.mu.RUnlock()
	m := DeepSearchInMap(flat.data, paths[:len(paths)-1]...)
	dd = m[paths[len(paths)-1]]
	flat.keyMap.Store(key, dd)
	return dd
}
