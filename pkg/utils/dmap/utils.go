package dmap

import (
	"fmt"
	"github.com/spf13/cast"
	"reflect"
)

// lookup Flattened multilevel map, add prefix + sep if there is one
// If no prefix, let the depth map be concatenated using sep
func lookup(prefix, sep string, sourceData, destData map[string]interface{}) {
	for key, value := range sourceData {
		fullIndex := fmt.Sprintf("%s%s%s", prefix, sep, key)
		if prefix == "" {
			fullIndex = fmt.Sprintf("%s", key)
		}
		if dd, err := cast.ToStringMapE(value); err != nil {
			//  The last level, can not be split
			destData[fullIndex] = value
		} else {
			// Have next level, go deep
			lookup(fullIndex, sep, dd, destData)
		}
	}
}

// DeepSearchInMap ...
func DeepSearchInMap(sourceMap map[string]interface{}, paths ...string) map[string]interface{} {
	// Deep Copy => Does not affect the original data
	newMap := make(map[string]interface{})
	for key, value := range sourceMap {
		newMap[key] = value
	}
	return deepSearch(newMap, paths)
}

// deepSearch Get data when the structure diagram is not flattened.
func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	var newMap = make(map[string]interface{})
	for _, searchKey := range path {
		subMap, ok := m[searchKey]
		// field to find value of searchKey
		if !ok {
			newMap = make(map[string]interface{})
			m[searchKey] = newMap
			m = newMap
			continue
		}
		newMap, ok = subMap.(map[string]interface{})
		if !ok {
			newMap = make(map[string]interface{})
			m[searchKey] = newMap
		}
		m = newMap
	}
	return m
}

// ToMapStringInterface convert map[interface{}]interface{} to map[string]interface{}
func ToMapStringInterface(src map[interface{}]interface{}) map[string]interface{} {
	dest := make(map[string]interface{})
	for key, value := range src {
		dest[fmt.Sprintf("%v", key)] = value
	}
	return dest
}

// MergeStringMap merge two map[string]interface to one Rule is 1. If dest has no value but src has => put src key value in desc 2. If src and dest both have key but value is different type => keep dest one 3. If src and dest both have key and is same type => using src one
func MergeStringMap(dest, src map[string]interface{}) {
	for sk, sv := range src {
		tv, ok := dest[sk]
		if !ok {
			dest[sk] = sv
			continue
		}
		svType := reflect.TypeOf(sv)
		tvType := reflect.TypeOf(tv)
		if svType != tvType {
			continue
		}
		switch ttv := tv.(type) {
		case map[interface{}]interface{}:
			tsv := sv.(map[interface{}]interface{})
			ssv := ToMapStringInterface(tsv)
			stv := ToMapStringInterface(ttv)
			MergeStringMap(stv, ssv)
			dest[sk] = stv
		case map[string]interface{}:
			MergeStringMap(ttv, sv.(map[string]interface{}))
			dest[sk] = ttv
		default:
			dest[sk] = sv
		}
	}
}
