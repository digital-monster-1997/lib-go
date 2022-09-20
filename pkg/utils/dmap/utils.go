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
