package dmap

import (
	"fmt"
	"github.com/DigitakMonster1997/lib-go/pkg/utils/dcast"
)

// lookup Flattened multilevel map, add prefix + sep if there is one
// If no prefix, let the depth map be concatenated using sep
func lookup(prefix, sep string, sourceData, destData map[string]interface{}) {
	for key, value := range sourceData {
		fullIndex := fmt.Sprintf("%s%s%s", prefix, sep, key)
		if prefix == "" {
			fullIndex = fmt.Sprintf("%s", key)
		}
		if dd, err := dcast.ToStringMapE(value); err != nil {
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
